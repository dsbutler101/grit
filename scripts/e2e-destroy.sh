#!/usr/bin/env bash

# Script for manually destroying e2e environments.
# Fallback in case the pipeline fails to destroy the resources.

# Setup process:
# 1. Find the terraform state to delete in GitLab, e.g. https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit/-/terraform
# 2. Configure input variables:
#    - TF_STATE_NAME: the name of the terraform state to destroy
#    - GITLAB_USER: your gitlab username
#    - GITLAB_TOKEN_TERRAFORM: a gitlab personal access token with API scope
# 3. Configure provider credentials for AWS/google/gitlab
# 4. Run the script

set -euo pipefail

# input variables
project_id="${PROJECT_ID:-48756626}"
gitlab_user="${GITLAB_USER:?Specify a GITLAB_USER to get terraform state}"
gitlab_token="${GITLAB_TOKEN_TERRAFORM:?Specify a GITLAB_TOKEN to get terraform state}"
state_name="${TF_STATE_NAME:?Specify a TF_STATE_NAME to destroy}"
timeout="${TIMEOUT:-20m}"

wkdir="$(mktemp -d)"
echo "Initializing terraform in ${wkdir}"

# configure terraform http backend
export TF_HTTP_ADDRESS="https://gitlab.com/api/v4/projects/${project_id}/terraform/state/${state_name}"
export TF_HTTP_LOCK_ADDRESS="${TF_HTTP_ADDRESS}/lock"
export TF_HTTP_UNLOCK_ADDRESS="${TF_HTTP_ADDRESS}/lock"
export TF_HTTP_USERNAME="${gitlab_user}"
export TF_HTTP_PASSWORD="${gitlab_token}"
export TF_HTTP_LOCK_METHOD=POST
export TF_HTTP_UNLOCK_METHOD=DELETE
export CI_PROJECT_ID="${project_id}"
# cache plugins here so subsequent runs don't have to download everything again
export TF_PLUGIN_CACHE_DIR="${PWD}/.terraform"
mkdir -p "${TF_PLUGIN_CACHE_DIR}"

# ensure we authenticate with the main providers so we can destroy resources
cat >"${wkdir}/main.tf" <<EOF
terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    google = {
      source = "hashicorp/google"
    }
    gitlab = {
      source = "gitlabhq/gitlab"
    }
  }

  backend "http" {}
}

provider "gitlab" {}
provider "aws" {}
provider "google" {}
EOF

terraform -chdir="${wkdir}" init

states="$(terraform -chdir="${wkdir}" state list)"
echo "Found remote states:"
echo "${states}"
# remove kubectl_ resources - we can't connect to clusters with only remote state
if echo "${states}" | grep -q kubectl_; then
  echo "Removing kubectl_ resources from remote state..."
  terraform -chdir="${wkdir}" state list | grep kubectl_ | xargs -d '\n' -I {} -r terraform -chdir="${wkdir}" state rm '{}'
fi
mage -t "${timeout}" terraform:initAndDestroy "${wkdir}"
