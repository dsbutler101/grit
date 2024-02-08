#!/usr/bin/env bash

set -eo pipefail

function decrypt_secret {
  local token=$${1}
  local secret=$${2}

  curl -s \
      -X POST \
      -H "Authorization: Bearer $${token}" \
      -H "Content-Type: application/json" \
      -d "{\"ciphertext\":\"$${secret}\"}" \
      'https://cloudkms.googleapis.com/v1/${kms_key}:decrypt' | jq -r .plaintext | base64 -d
}

#
# Check required software existence
#
which gitlab-runner
gitlab-runner --version

#
# Install fleeting plugin
#

%{ if use_autoscaling }
curl -Lo /usr/local/bin/fleeting-plugin-googlecompute \
    "https://gitlab.com/gitlab-org/fleeting/fleeting-plugin-googlecompute/-/releases/${fleeting_googlecompute_plugin_version}/downloads/fleeting-plugin-googlecompute-linux-amd64"
chmod +x /usr/local/bin/fleeting-plugin-googlecompute
%{ endif }

#
# Retrieve runner token from the encrypted secret
#
apk add -U jq

TOKEN=$(curl --silent --header "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token | jq -r .access_token)

# Prepare runner token from the secret
RUNNER_TOKEN=$(decrypt_secret "$${TOKEN}" "${runner_token}")

# Prepare SSH private key from the secret
decrypt_secret "$${TOKEN}" "${runner_ssh_key}" > /etc/gitlab-runner/ssh-key.priv
chmod 0600 /etc/gitlab-runner/ssh-key.priv

#
# Register runner with the authentication token decrypted from the secret
#
register_flags=()
register_flags+=(--config "/etc/gitlab-runner/config.toml")
register_flags+=(--name "${name}")
register_flags+=(--url "${gitlab_url}")
register_flags+=(--token "$${RUNNER_TOKEN}")
register_flags+=(--template-config "/etc/gitlab-runner/config-template.toml")

# shellcheck disable=SC2068
gitlab-runner register --non-interactive $${register_flags[@]}

#
# Start runner
#
exec /entrypoint "$@"
