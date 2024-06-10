#!/usr/bin/env bash

set -e
set -u
set -o pipefail

: "${REPO_BASE_URL:=https://gitlab.com/gitlab-org/gl-openshift/gitlab-runner-operator}"
: "${ARTIFACT:=operator.k8s.yaml}"
: "${DEST_FILE_NAME:=manifests.yaml}"
: "${DEST_DIR:=$( cd "$(dirname "${BASH_SOURCE[0]}")" && pwd )/versions}"
: "${UPDATE_CURRENT:=1}"

main() {
  local tag name current=''

  while IFS=$'\t' read -r tag name ; do
    downloadVersion "$tag" "$name"

    [[ $UPDATE_CURRENT == 1 && -z "$current" ]] && {
      echo >&2 "## setting current to ${tag}"
      current="${tag}"
      ln -sf -T "${tag}" "${DEST_DIR}/current"
    }
  done < <( getVersions | sort )
}

curl() {
  command curl -LfSs "$@"
}

getVersions() {
  curl -H 'Accept: application/json' "${REPO_BASE_URL}/-/releases" \
    | jq -r '.[] | [.tag, .name] | @tsv'
}

sort() {
  command sort -Vr
}

downloadVersion() {
  local tag="$1"
  local name="$2"
  local src="${REPO_BASE_URL}/-/releases/${name}/downloads/${ARTIFACT}"
  local dest="${DEST_DIR}/${tag}/${DEST_FILE_NAME}"

  mkdir -p "$( dirname "${dest}" )"

  echo >&2 "## downloading ${ARTIFACT} from release ${name} (tag: ${tag})"
  curl "${src}" -o "${dest}"
}

main "$@"
