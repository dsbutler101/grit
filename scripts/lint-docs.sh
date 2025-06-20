#!/usr/bin/env bash

set -eo pipefail

GIT_ROOT=$(cd "${BASH_SOURCE%/*}" && git rev-parse --show-toplevel)
VALE_MIN_ALERT_LEVEL=${VALE_MIN_ALERT_LEVEL:-}
ERROR_RESULTS=0

echo "Lint prose"
if command -v vale >/dev/null 2>&1; then
    args=()
    if [ -n "${VALE_MIN_ALERT_LEVEL}" ]; then
        args+=("--minAlertLevel" "${VALE_MIN_ALERT_LEVEL}")
    fi
    vale --config "${GIT_ROOT}/.vale.ini" --glob "!**/.terraform/**" "${args[@]}" "${GIT_ROOT}" || ((ERROR_RESULTS++))
else
    echo "vale is missing, please install it from https://vale.sh/docs/vale-cli/installation/"
fi

echo "Lint Markdown"
if command -v markdownlint-cli2 >/dev/null 2>&1; then
    markdownlint-cli2 "**/*.md" '!**/.terraform/**' '!CHANGELOG.md' || ((ERROR_RESULTS++))
else
    echo "markdownlint-cli2 is missing, please install it from https://github.com/DavidAnson/markdownlint-cli2#install"
fi

if [ "${ERROR_RESULTS}" -ne 0 ]; then
    echo "✖ ${ERROR_RESULTS} lint test(s) failed. Review the log carefully to see full listing."
    exit 1
else
    echo "✔ Linting passed"
    exit 0
fi
