module "runner" {
  source = "../internal/"

  url             = var.gitlab.url
  token           = var.gitlab.runner_token
  namespace       = var.namespace
  name            = coalesce(var.name_override, var.metadata.name)
  config_template = var.config_template
}
