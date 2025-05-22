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
        latestVersion {
          id
        }
      }
    }
  }
}"

# shellcheck disable=SC2016
data="$(gojq -nc --arg q "${query}" '{"query": $q}')"

resp="$(curl -s "https://gitlab.com/api/graphql?private_token=${GITLAB_TOKEN_TERRAFORM}" -H "Content-Type: application/json" -X POST -d "${data}")"

# Get all old states (created more than 24 hours ago with u[0-9] naming pattern)
old_states="$(echo "${resp}" | gojq -r '.data.project.terraformStates.nodes[] | select((.name | test("^u[0-9]")) and (now - (.createdAt | fromdate) > 86400))')"

# Filter for states to destroy: old, not locked, with latestVersion
states_to_destroy="$(echo "${old_states}" | gojq -r 'select((.lockedAt == null) and (.latestVersion != null)) | .name')"

declare -a failedStates

# destroy
for state in ${states_to_destroy}; do
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

# Find states with no latestVersion (null)
states_with_no_version="$(echo "${old_states}" | gojq -r 'select(.latestVersion == null) | .name')"
if [[ -n "${states_with_no_version}" ]]; then
  echo -e "\n---\n"
  echo "Following E2E TF states have no stored state (no latestVersion)."
  echo -e "They will require manual cleanup.\n"

  for state in ${states_with_no_version}; do
    echo "- ${state}"
  done
fi

# Find locked states
locked_states="$(echo "${old_states}" | gojq -r 'select(.lockedAt != null) | .name')"
if [[ -n "${locked_states}" ]]; then
  echo -e "\n---\n"
  echo "Following E2E TF states are older than 24 hours and are currently locked."
  echo -e "They will require manual cleanup.\n"

  for state in ${locked_states}; do
    echo "- ${state}"
  done
fi
