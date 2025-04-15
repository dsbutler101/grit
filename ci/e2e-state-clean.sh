#!/usr/bin/env bash

set -eo pipefail

# graphql query
query="{
  project(fullPath: \"gitlab-org/ci-cd/runner-tools/grit\") {
    terraformStates(first: 100) {
      nodes {
        name
        createdAt
        lockedAt
      }
    }
  }
}"

# shellcheck disable=SC2016
data="$(gojq -nc --arg q "${query}" '{"query": $q}')"

resp="$(curl -s "https://gitlab.com/api/graphql?private_token=${GITLAB_TOKEN_TERRAFORM}" -H "Content-Type: application/json" -X POST -d "${data}")"

# filter out non-e2e states and states those created in last 24 hours that are not locked
states="$(echo "${resp}" | gojq -r '.data.project.terraformStates.nodes[] | select((.name | test("^u[0-9]")) and (now - (.createdAt | fromdate) > 86400) and (.lockedAt == null)) | .name')"

declare -a failedStates

# destroy
for state in ${states}; do
  echo -e "\033[31;1mDestroying terraform state \033[37;1m${state}\033[0m"
  export TF_STATE_NAME="${state}"
  ./scripts/e2e-destroy.sh || failedStates+=("${state}")
done

if [[ ${#failedStates[@]} -gt 0 ]]; then
  echo -e "\n---\n"
  echo "Following E2E TF states were failed to destroy."
  echo -e "They will require manual cleanup.\n"

  for failedState in "${failedStates[@]}"; do
    echo "- ${failedState}"
  done
fi

states="$(echo "${resp}" | gojq -r '.data.project.terraformStates.nodes[] | select((.name | test("^u[0-9]")) and (now - (.createdAt | fromdate) > 86400) and (.lockedAt != null)) | .name')"
if [[ -n "${states}" ]]; then
  echo -e "\n---\n"
  echo "Following E2E TF states are older than 24 hours and are currently locked."
  echo -e "They will require manual cleanup.\n"

  for state in ${states}; do
    echo "- ${state}"
  done
fi
