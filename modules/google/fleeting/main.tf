#######################
# METADATA VALIDATION #
#######################

module "validate_name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate_support" {
  source   = "../../internal/validation/support"
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

  source = "./gce"

  name   = var.metadata.name
  labels = var.metadata.labels

  google_project        = var.google_project
  subnetwork_project    = var.subnetwork_project
  access_config_enabled = var.access_config_enabled
  google_zone           = var.google_zone
  service_account_email = var.service_account_email

  machine_type = var.machine_type
  disk_type    = var.disk_type
  disk_size_gb = var.disk_size_gb
  source_image = var.source_image

  vpc                                     = var.vpc
  manager_subnet_cidr                     = var.manager_subnet_cidr
  additional_tags                         = var.additional_tags
  cross_vm_deny_egress_destination_ranges = var.cross_vm_deny_egress_destination_ranges
}
