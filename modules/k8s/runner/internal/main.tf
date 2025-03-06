locals {
  config_template_name = format("%s-%s", var.name, "config-template")
  envvars_name         = format("%s-%s", var.name, "envvars")
  manifest = yamlencode({
    apiVersion = "apps.gitlab.com/v1beta2"
    kind       = "Runner"
    metadata = {
      name      = var.name
      namespace = var.namespace
    }
    spec = merge({
      gitlabUrl     = var.url
      token         = var.name
      locked        = var.locked
      protected     = var.protected
      tags          = join(",", var.runner_tags)
      runUntagged   = length(var.runner_tags) == 0 ? true : var.run_untagged
      interval      = var.check_interval
      concurrent    = var.concurrent
      config        = var.config_template == "" ? null : local.config_template_name
      env           = length(var.envvars) == 0 ? null : local.envvars_name
      podSpec       = var.pod_spec_patches
      runnerImage   = var.runner_image
      helperImage   = var.helper_image
      logLevel      = var.log_level
      listenAddress = var.listen_address
    }, var.runner_opts)
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
  envvars = yamlencode({
    apiVersion = "v1"
    kind       = "ConfigMap"
    metadata = {
      name      = local.envvars_name
      namespace = var.namespace
    }
    data = var.envvars
  })

  config_template_check = length(var.config_template) == 0 || (length(var.config_template) > 0 && strcontains(var.config_template, "[[runners]]"))
}

module "check_config_template" {
  source  = "../../../internal/validation/fail_validation"
  message = local.config_template_check ? "" : "The config template must contain the definition of [[runners]]."
}

resource "terraform_data" "token_secret" {
  input = local.token_secret
}

resource "terraform_data" "config_template" {
  input = local.config_template
}

resource "kubectl_manifest" "token_secret" {
  yaml_body = terraform_data.token_secret.input
  wait      = true
  force_new = true
}

resource "kubectl_manifest" "config_template" {
  count     = var.config_template == "" ? 0 : 1
  yaml_body = local.config_template
  wait      = true
  force_new = true
}

resource "kubectl_manifest" "envvars" {
  count     = length(var.envvars) == 0 ? 0 : 1
  yaml_body = local.envvars
}

resource "kubectl_manifest" "manifest" {
  yaml_body = local.manifest
  wait      = true
  force_new = true

  depends_on = [
    kubectl_manifest.token_secret,
    kubectl_manifest.config_template,
    kubectl_manifest.envvars,
  ]

  lifecycle {
    replace_triggered_by = [
      terraform_data.token_secret,
      terraform_data.config_template,
    ]
  }
}
