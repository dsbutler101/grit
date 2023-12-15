###############
# TEST MODULE #
###############

module "test-module" {
  source = "./internal"

  gitlab_project_id         = var.gitlab_project_id
  gitlab_url                = var.gitlab_url
  gitlab_runner_description = var.gitlab_runner_description
  gitlab_runner_tags        = var.gitlab_runner_tags
  runner_token              = var.runner_token

  manager_service       = var.manager_service
  fleeting_service      = var.fleeting_service
  fleeting_os           = var.fleeting_os
  executor              = var.executor
  scale_min             = var.scale_min
  scale_max             = var.scale_max
  idle_percentage       = var.idle_percentage
  capacity_per_instance = var.capacity_per_instance

  ami                                  = var.ami
  instance_type                        = var.instance_type
  asg_storage_type                     = var.asg_storage_type
  asg_storage_size                     = var.asg_storage_size
  asg_storage_throughput               = var.asg_storage_throughput
  aws_zone                             = var.aws_zone
  aws_vpc_cidr                         = var.aws_vpc_cidr
  aws_vpc_subnet_cidr                  = var.aws_vpc_subnet_cidr
  aws_vpc_id                           = var.aws_vpc_id
  aws_vpc_subnet_id                    = var.aws_vpc_subnet_id
  macos_cores_per_license              = var.macos_cores_per_license
  macos_required_license_count_per_asg = var.macos_required_license_count_per_asg

  name   = var.name
  labels = var.labels
}

check "vpc" {

  // Only cidrs or ids should be set
  assert {
    condition     = (var.aws_vpc_cidr != "" && var.aws_vpc_subnet_cidr != "") || (var.aws_vpc_id != "" && var.aws_vpc_subnet_id != "")
    error_message = "Only set CIDRs or IDs for VPC configuration. Only aws_vpc_cidr and aws_vpc_subnet_cidr should be set or aws_vpc_id and aws_vpc_subnet_id should be set."
  }
}
