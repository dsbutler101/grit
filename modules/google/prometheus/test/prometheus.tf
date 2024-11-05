#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

##########################
# PROMETHEUS TEST CONFIG #
##########################

module "prometheus" {
  source = "../internal"

  name   = var.metadata.name
  labels = var.metadata.labels

  google_project = var.google_project
  google_zone    = var.google_zone

  service_account_email = var.service_account_email

  machine_type = var.machine_type
  boot_disk    = var.boot_disk
  data_disk    = var.data_disk

  prometheus_version = var.prometheus_version

  node_exporter_version = var.node_exporter_version
  node_exporter_port    = var.node_exporter_port

  prometheus_external_labels = var.prometheus_external_labels
  mimir                      = var.mimir
  runner_manager_nodes       = var.runner_manager_nodes

  vpc = var.vpc
}
