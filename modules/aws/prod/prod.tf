###############
# PROD MODULE #
###############

locals {
  use_case = "${var.fleeting_service}-${var.manager_provider}-${var.executor}"
  use_case_maturity = tomap({
    "ec2-ec2-docker-autoscaler" = "alpha"
  })
  maturity = try(local.use_case_maturity[local.use_case], "unsupported")
}

module "prod-module" {
  source = "./internal"

  fleeting_service = var.fleeting_service
  fleeting_os      = var.fleeting_os
  ami              = var.ami
  executor         = var.executor
  instance_type    = var.instance_type
  aws_vpc_cidr     = var.aws_vpc_cidr

  capacity_per_instance = var.capacity_per_instance
  scale_min             = var.scale_min
  scale_max             = var.scale_max

  manager_provider = var.manager_provider

  gitlab_project_id         = var.gitlab_project_id
  gitlab_runner_description = var.gitlab_runner_description
  gitlab_runner_tags        = var.gitlab_runner_tags
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

  // min_maturity beta can be satisfied by only stable
  assert {
    condition     = var.min_maturity != "stable" || local.maturity == "stable"
    error_message = "Maturity is ${local.maturity} but min_maturity is ${var.min_maturity}"
  }
}