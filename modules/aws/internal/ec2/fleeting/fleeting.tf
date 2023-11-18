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

  asg_storage_type       = "gp3"
  asg_storage_throughput = 750
  asg_ami_id             = var.ami
  asg_instance_type      = var.instance_type
  asg_subnet_cidr        = "10.0.0.0/24"
  aws_vpc_cidr           = "10.0.0.0/24"
  name                   = var.name
}

module "macos" {
  count  = var.fleeting_os == "macos" ? 1 : 0
  source = "./macos"

  asg_ami_id        = var.ami
  asg_instance_type = var.instance_type
  asg_subnet_cidr   = "10.0.0.0/24"

  aws_vpc_cidr = "10.0.0.0/24"
  name         = var.name
}
