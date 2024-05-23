#!/usr/bin/env bash

#+ This pulls down the upstream gitlab runner operator manifests, massages them
#+ into a format that we can pass on to terraform, and puts them into a file
#+ that will be automatically picked up by terraform.
#+
#+ Note:
#+   The output file needs to be named '.tf.json', otherwise terraform will ignore it

set -e
set -u
set -o pipefail

: "${REPO_BASE_URL:=https://gitlab.com/gitlab-org/gl-openshift/gitlab-runner-operator}"
: "${REPO_PATH:=operator.k8s.yaml}"
: "${REV:=master}"
: "${DEST_FILE:=locals_generated.tf.json}"

main() {
  local rawURL="${REPO_BASE_URL}/-/raw/${REV}/${REPO_PATH}"
  local dir="$( cd "$(dirname "${BASH_SOURCE[0]}")" &&  pwd )"
  local destFull="${dir}/${DEST_FILE}"

  curl -fLsS "${rawURL}" \
    | yq -se '
        {
          "//": "auto generated, do not edit",
          "locals": {
            "operator_manifests": (
              map(
                "\(.apiVersion)::\(.kind)::\(.metadata.namespace//"_cluster_scoped_")::\(.metadata.name)" as $name
                | { $name : . }
              ) | add
            )
          }
        }
      ' \
    > "${dir}/${DEST_FILE}"

  echo >&2 "# Updated '${dir}/${DEST_FILE}' with restructured content from '${rawURL}'"
}

main "$@"
