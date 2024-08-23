module "runner" {
  source = "../internal/"

  url             = var.gitlab.url
  token           = var.gitlab.runner_token
  namespace       = var.namespace
  name            = coalesce(var.name_override, var.metadata.name)
  concurrent      = var.concurrent
  check_interval  = var.check_interval
  locked          = var.locked
  protected       = var.protected
  run_untagged    = var.run_untagged
  runner_tags     = var.runner_tags
  config_template = var.config_template
}
