#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

########################
# FLEETING TEST MODULE #
########################

module "ec2" {
  count  = var.service == "ec2" ? 1 : 0
  source = "../internal/ec2"

  vpc = var.vpc

  name                        = var.metadata.name
  os                          = var.os
  ami                         = var.ami
  instance_type               = var.instance_type
  scale_min                   = var.scale_min
  scale_max                   = var.scale_max
  labels                      = var.metadata.labels
  storage_type                = var.storage_type
  storage_size                = var.storage_size
  storage_throughput          = var.storage_throughput
  macos_license_count_per_asg = var.macos_license_count_per_asg
  macos_cores_per_license     = var.macos_cores_per_license
  security_group_ids          = var.security_group_ids
  install_cloudwatch_agent    = var.install_cloudwatch_agent
  cloudwatch_agent_json       = var.cloudwatch_agent_json
  instance_role_profile_name  = var.instance_role_profile_name
}
