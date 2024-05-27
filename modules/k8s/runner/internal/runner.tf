locals {
  manifest = yamlencode({
    apiVersion = "apps.gitlab.com/v1beta2"
    kind       = "Runner"
    metadata = {
      name      = var.name
      namespace = var.namespace
    }
    spec = {
      gitlabUrl = var.url
      token     = var.name
      locked    = true
    }
  })
  token_secret = yamlencode({
    apiVersion = "v1"
    kind       = "Secret"
    metadata = {
      name      = var.name
      namespace = var.namespace
    }
    data = {
      runner-registration-token = base64encode(var.token)
    }
  })
}

resource "kubectl_manifest" "token_secret" {
  yaml_body = local.token_secret
}
resource "kubectl_manifest" "manifest" {
  yaml_body = local.manifest
  depends_on = [
    kubectl_manifest.token_secret
  ]
}
