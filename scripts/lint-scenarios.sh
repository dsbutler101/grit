#!/usr/bin/env bash

# Checks some aspects of scenario docs
# See the comments on the `assert::` functions for more details

set -e
set -u
set -o pipefail

readonly BASE_DIR="$( cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd )"
readonly DOCS_DIR="${BASE_DIR}/docs/scenarios"
readonly DOC_INDEX="${DOCS_DIR}/index.md"
readonly SCENARIO_DIR="${BASE_DIR}/scenarios"

main() {
  local name
  local rc=0

  while read -r -d $'\0' name ; do
    if validateScenario "$name" ; then
      log "${name}" '✔'
    else
      rc=$(( rc | $? ))
    fi
  done < <( find "${SCENARIO_DIR}" -type d -links 2 -printf '%P\0' )

  return $rc
}

validateScenario() {
  local name="$1"
  local rc=0

  assert::docExists "$name" || rc=$(( rc | $? ))
  assert::docIsLinked "$name" || rc=$(( rc | $? ))
  assert::docTfInputs "$name" || rc=$(( rc | $? ))

  return $rc
}

log() {
  local name="$1"
  local msg="$2"
  shift 2

  local data
  {
    echo "[${name}] ${msg}"

    for data ; do
      sed 's/^/\t| /g' <<< "$data"
    done
  } >&2
}

fail() {
  log "$@"
  return 1
}

# Checks, if a scenario has documentation
assert::docExists() {
  local name="$1"
  local expectedFiles=(
    "${DOCS_DIR}/${name}/index.md"
  )
  local f rc=0

  for f in "${expectedFiles[@]}" ; do
    test -r "$f" || {
      fail "$name" "expected doc file '$f' does not exist" || rc=$(( rc | $? ))
    }
  done

  return $rc
}

# Checks, if a individual scenario's doc is linked in the scenario doc index
assert::docIsLinked() {
  local name="$1"

  grep -qF "$name" "${DOC_INDEX}" || {
    fail "$name" "expected docs to be linked in index"
  }
}

# Checks, if all terraform inputs are doc'ed
assert::docTfInputs() {
  local -a tfVars docVars
  local diffBlob

  mapfile -t tfVars < <(
    grep -RIPoh 'variable\s+"\K[^"]+(?=")' "${SCENARIO_DIR}/${name}" \
      | sort
  )

  mapfile -t docVars < <(
    sed -n '/begin: input vars/,/end: input vars/p' "${DOCS_DIR}/${name}/index.md" \
      | grep -Po '^\|\s+`\K[^`]+(?=`)' \
      | sort
  )

  diffBlob="$(
    diff \
      --label '~variables.tf' \
      --label '~doc.md' \
      -Naur \
      <( array::join $'\n' "${tfVars[@]}" ) \
      <( array::join $'\n' "${docVars[@]}" )
  )" || {
    fail "$name" "set of terraform input vars does not match the set of documented vars:" "${diffBlob}"
  }
}

array::join() {
  local IFS="$1" ; shift
  echo "$*"
}


main "$@"
