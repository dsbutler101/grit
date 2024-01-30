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

  name                                   = var.metadata.name
  os                                     = var.os
  ami                                    = var.ami
  instance_type                          = var.instance_type
  scale_min                              = var.scale_min
  scale_max                              = var.scale_max
  idle_percentage                        = var.idle_percentage
  labels                                 = var.metadata.labels
  storage_type                           = var.storage_type
  storage_size                           = var.storage_size
  storage_throughput                     = var.storage_throughput
  macos_required_license_counter_per_asg = var.macos_required_license_counter_per_asg
  macos_cores_per_license                = var.macos_cores_per_license
}
