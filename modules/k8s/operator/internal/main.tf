locals {
  operator_project_id = "22848448"

  manifest = startswith(var.override_manifests, "file://") ? file(trimprefix(var.override_manifests, "file://")) : data.http.manifest.response_body

  # It's guaranteed that the Namespace will be the first resource in the manifest, we match by its name
  namespace = try(regex("\\s+name:\\s+([^\\s]+)", local.manifest)[0], "default")

  operator_version = var.operator_version == "latest" ? module.latest_operator_version.tags[0] : var.operator_version

  download_url = "https://gitlab.com/api/v4/projects/${local.operator_project_id}/packages/generic/gitlab-runner-operator/${local.operator_version}/manifests/operator.k8s.yaml"
}

data "http" "manifest" {
  url = can(regex("^https?://", var.override_manifests)) ? var.override_manifests : local.download_url
}

module "latest_operator_version" {
  source = "../../../../modules/internal/gitlab_tags"

  project_id = local.operator_project_id
}

resource "kubectl_manifest" "operator_resources" {
  for_each = {
    for resource in [
      for yaml in split("\n---\n", local.manifest) : {
        decoded = yamldecode(yaml),
        encoded = yaml
      }
    ] : "${resource.decoded.kind}:${resource.decoded.metadata.name}" => resource.encoded
  }
  yaml_body        = each.value
  wait             = true
  apply_only       = true
  force_new        = true
  wait_for_rollout = false
}
