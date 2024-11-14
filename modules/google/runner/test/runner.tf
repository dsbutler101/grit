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
  google_zone    = var.google_zone

  node_exporter = var.node_exporter

  service_account_email = var.service_account_email
  machine_type          = var.machine_type
  disk_type             = var.disk_type
  disk_size_gb          = var.disk_size_gb
  runner_version        = var.runner_version

  concurrent     = var.concurrent
  check_interval = var.check_interval
  log_level      = var.log_level

  runner_metrics_listener = var.runner_metrics_listener
  listen_address          = var.listen_address

  gitlab_url          = var.gitlab_url
  runner_token        = var.runner_token
  request_concurrency = var.request_concurrency
  executor            = var.executor

  cache_gcs_bucket = var.cache_gcs_bucket

  runners_global_section = var.runners_global_section
  runners_docker_section = var.runners_docker_section
  default_docker_image   = var.default_docker_image

  fleeting_googlecompute_plugin_version = var.fleeting_googlecompute_plugin_version
  fleeting_instance_group_name          = var.fleeting_instance_group_name

  capacity_per_instance = var.capacity_per_instance
  max_use_count         = var.max_use_count
  max_instances         = var.max_instances
  autoscaling_policies  = var.autoscaling_policies

  runner_manager_additional_firewall_rules = var.runner_manager_additional_firewall_rules
  vpc                                      = var.vpc
}
