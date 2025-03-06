module "runner" {
  source = "../internal/"

  url              = var.gitlab.url
  token            = var.gitlab.runner_token
  namespace        = var.namespace
  name             = coalesce(var.name_override, var.metadata.name)
  concurrent       = var.concurrent
  check_interval   = var.check_interval
  locked           = var.locked
  protected        = var.protected
  run_untagged     = var.run_untagged
  runner_tags      = var.runner_tags
  config_template  = var.config_template
  envvars          = var.envvars
  pod_spec_patches = var.pod_spec_patches
  runner_image     = var.runner_image
  helper_image     = var.helper_image
  log_level        = var.log_level
  listen_address   = var.listen_address
  runner_opts      = var.runner_opts
}
