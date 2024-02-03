#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../../../internal/validation/support"
  use_case = var.service
  use_case_support = tomap({
    "ec2" = "experimental"
  })
  min_support = var.metadata.min_support
}

########################
# FLEETING PROD MODULE #
########################

module "ec2" {
  count  = var.service == "ec2" ? 1 : 0
  source = "../internal/ec2"

  vpc = var.vpc

  name                        = var.metadata.name
  os                          = var.os
  ami                         = var.ami
  instance_type               = var.instance_type
  labels                      = var.metadata.labels
  storage_type                = var.storage_type
  storage_size                = var.storage_size
  storage_throughput          = var.storage_throughput
  macos_license_count_per_asg = var.macos_license_count_per_asg
  macos_cores_per_license     = var.macos_cores_per_license
  scale_min                   = var.scale_min
  scale_max                   = var.scale_max
}
