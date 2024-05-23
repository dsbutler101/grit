terraform {
  required_version = ">= 0.13"

  required_providers {
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
}

resource "kubectl_manifest" "operator_resources" {
  for_each = {
    for n, m in local.operator_manifests : n => yamlencode(m)
  }
  yaml_body = each.value
}
