#!/usr/bin/env bash

set -e
set -u
set -o pipefail

: "${REPO_BASE_URL:=https://gitlab.com/gitlab-org/gl-openshift/gitlab-runner-operator}"
: "${REPO_PATH:=operator.k8s.yaml}"
: "${REV:=master}"
: "${DEST_FILE:=manifests.yaml}"

main() {
  local rawURL="${REPO_BASE_URL}/-/raw/${REV}/${REPO_PATH}"
  local dir="$( cd "$(dirname "${BASH_SOURCE[0]}")" &&  pwd )"
  local destFull="${dir}/${DEST_FILE}"

  curl -fLsS "${rawURL}" -o "${destFull}"

  echo >&2 "# Updated '${destFull}' (source: '${rawURL}')"
}

main "$@"
