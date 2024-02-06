#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

######################
# RUNNER TEST CONFIG #
######################

module "runner" {
  source = "../internal"

  name   = var.metadata.name
  labels = var.metadata.labels

  google_project = var.google_project

  service_account_email = var.service_account_email
  machine_type          = var.machine_type
  disk_type             = var.disk_type
  disk_size_gb          = var.disk_size_gb
  runner_version        = var.runner_version

  concurrent     = var.concurrent
  check_interval = var.check_interval
  log_level      = var.log_level
  listen_address = var.listen_address

  gitlab_url          = var.gitlab_url
  runner_token        = var.runner_token
  request_concurrency = var.request_concurrency

  cache_gcs_bucket = var.cache_gcs_bucket

  runners_global_section = var.runners_global_section
  runners_docker_section = var.runners_docker_section

  vpc = var.vpc
}
