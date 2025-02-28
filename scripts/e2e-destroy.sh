#!/usr/bin/env bash

# Script for manually destroying e2e environments.
# Fallback in case the pipeline fails to destroy the resources.
# Assumes an up-to-date working tree.

# Setup process:
# 1. Find the terraform state to delete in GitLab, e.g. https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit/-/terraform
# 2. Configure input variables:
#    - TF_STATE_NAME: the name of the terraform state to destroy
#    - COMMIT_SHA: the commit SHA to checkout, found by looking up the pipeline ID in the state name.
#    - MODULE_PATH: the relative path to the terraform module, e.g. e2e/google
#    - GITLAB_USER: your gitlab username
#    - GITLAB_TOKEN_TERRAFORM: a gitlab personal access token with API scope
# 3. Configure required `TF_VAR` variables to make module destroy work
# 4. Run the script and probably go back to 3.

set -euo pipefail

# input variables
project_id="${PROJECT_ID:-48756626}"
gitlab_user="${GITLAB_USER:?Specify a GITLAB_USER to get terraform state}"
gitlab_token="${GITLAB_TOKEN_TERRAFORM:?Specify a GITLAB_TOKEN to get terraform state}"
state_name="${TF_STATE_NAME:?Specify a TF_STATE_NAME to destroy}"
commit_sha="${COMMIT_SHA:?Specify a GRIT COMMIT_SHA to clone}"
path="${MODULE_PATH:?Specify a relative MODULE_PATH to the terraform module}"
timeout="${TIMEOUT:-20m}"

wkdir="$(mktemp -d)"
echo "Checking out ${commit_sha} to ${wkdir}"

git worktree add "$wkdir" "$commit_sha"

# configure terraform http backend
export TF_HTTP_ADDRESS="https://gitlab.com/api/v4/projects/${project_id}/terraform/state/${state_name}"
export TF_HTTP_LOCK_ADDRESS="${TF_HTTP_ADDRESS}/lock"
export TF_HTTP_UNLOCK_ADDRESS="${TF_HTTP_ADDRESS}/lock"
export TF_HTTP_USERNAME="${gitlab_user}"
export TF_HTTP_PASSWORD="${gitlab_token}"
export TF_HTTP_LOCK_METHOD=POST
export TF_HTTP_UNLOCK_METHOD=DELETE
export CI_PROJECT_ID="${project_id}"

trap "git worktree remove $wkdir" EXIT

mage terraform:initAndDestroy "$wkdir/$path" "${timeout}"
