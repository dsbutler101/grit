#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../../../internal/validation/support"
  use_case = "gce"
  use_case_support = tomap({
    "gce" = "experimental"
  })
  min_support = var.metadata.min_support
}

########################
# FLEETING PROD MODULE #
########################

module "gce" {
  count = var.fleeting_service == "gce" ? 1 : 0

  source = "../internal/gce"

  name   = var.metadata.name
  labels = var.metadata.labels

  google_project = var.google_project

  service_account_email = var.service_account_email

  machine_type = var.machine_type
  disk_type    = var.disk_type
  disk_size_gb = var.disk_size_gb
  source_image = var.source_image

  vpc                 = var.vpc
  manager_subnet_cidr = var.manager_subnet_cidr
}
