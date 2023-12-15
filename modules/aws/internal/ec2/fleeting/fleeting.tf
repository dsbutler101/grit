###############################
# INTERNAL EC2 INSTANCE GROUP #
###############################

module "common" {
  source = "./common"

  license_arn                      = try(module.macos.license-config-arn, "")
  jobs-host-resource-group-outputs = try(module.macos.jobs-host-resource-group-outputs, {})

  scale_min       = var.scale_min
  scale_max       = var.scale_max
  idle_percentage = var.idle_percentage

  asg_storage_type       = var.asg_storage_type
  asg_storage_size       = var.asg_storage_size
  asg_storage_throughput = var.asg_storage_throughput
  asg_ami_id             = var.ami
  asg_instance_type      = var.instance_type
  vpc_id                 = var.vpc_id
  subnet_id              = var.subnet_id

  labels = var.labels
  name   = var.name
}

module "macos" {
  count  = var.fleeting_os == "macos" ? 1 : 0
  source = "./macos"

  required_license_count_per_asg = var.macos_required_license_count_per_asg
  cores_per_license              = var.macos_cores_per_license

  labels = var.labels
  name   = var.name
}
