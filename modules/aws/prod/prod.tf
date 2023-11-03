###############
# PROD MODULE #
###############

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