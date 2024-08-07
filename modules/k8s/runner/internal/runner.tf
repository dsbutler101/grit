locals {
  config_template_name = format("%s-%s","config-template",var.name)
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
      config    = local.config_template_name
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
      runner-token = base64encode(var.token)
    }
  })
  config_template = yamlencode({
    apiVersion = "v1"
    kind       = "ConfigMap"
    metadata = {
      name      = local.config_template_name
      namespace = var.namespace
    }
    data = {
      "config.toml" = var.config_template
    }
  })
}

resource "kubectl_manifest" "token_secret" {
  yaml_body = local.token_secret
}

resource "kubectl_manifest" "config_template" {
  yaml_body = local.config_template
}

resource "kubectl_manifest" "manifest" {
  yaml_body = local.manifest
  depends_on = [
    kubectl_manifest.token_secret,
    kubectl_manifest.config_template
  ]
}
