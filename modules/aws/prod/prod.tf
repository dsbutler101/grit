###############
# PROD MODULE #
###############

locals {
  use_case = "${var.fleeting_service}-${var.manager_service}-${var.executor}"
  use_case_maturity = tomap({
    "ec2-ec2-docker-autoscaler" = "alpha"
  })
  maturity = try(local.use_case_maturity[local.use_case], "unsupported")
}

module "prod-module" {
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

check "maturity" {

  // min_maturity alpha can be satisfied by alpha, beta or stable
  assert {
    condition     = var.min_maturity != "alpha" || local.maturity == "alpha" || local.maturity == "beta" || local.maturity == "stable"
    error_message = "Maturity is ${local.maturity} but min_maturity is ${var.min_maturity}"
  }

  // min_maturity beta can be satisfied by beta or stable
  assert {
    condition     = var.min_maturity != "beta" || local.maturity == "beta" || local.maturity == "stable"
    error_message = "Maturity is ${local.maturity} but min_maturity is ${var.min_maturity}"
  }

  // min_maturity stable can be satisfied by only stable
  assert {
    condition     = var.min_maturity != "stable" || local.maturity == "stable"
    error_message = "Maturity is ${local.maturity} but min_maturity is ${var.min_maturity}"
  }
}

check "vpc" {

  // Only cidrs or ids should be set
  assert {
    condition     = (var.aws_vpc_cidr != "" && var.aws_vpc_subnet_cidr != "") || (var.aws_vpc_id != "" && var.aws_vpc_subnet_id != "")
    error_message = "Only set CIDRs or IDs for VPC configuration. Only aws_vpc_cidr and aws_vpc_subnet_cidr should be set or aws_vpc_id and aws_vpc_subnet_id should be set."
  }
}
