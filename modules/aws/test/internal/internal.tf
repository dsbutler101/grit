locals {
  access_key_id     = try(module.ec2-instance-group[0].fleeting_service_account_access_key_id, "")
  secret_access_key = try(module.ec2-instance-group[0].fleeting_service_account_secret_access_key, "")
  ssh_key_pem       = try(module.ec2-instance-group[0].ssh_key_pem, "")
  ssh_key_pem_name  = try(module.ec2-instance-group[0].ssh_key_pem_name, "")
  aws_asg_name      = try(module.ec2-instance-group[0].autoscaling_group_name, "")
  vpc_id            = try(module.vpc[0].vpc_id, var.aws_vpc_id)
  subnet_id         = try(module.vpc[0].subnet_id, var.aws_vpc_subnet_id)
}

#######################
# GITLAB REGISTRATION #
#######################

module "gitlab" {
  count  = var.runner_token == "" && var.gitlab_project_id != "" ? 1 : 0
  source = "../../../gitlab/internal"

  gitlab_project_id         = var.gitlab_project_id
  gitlab_runner_description = var.gitlab_runner_description
  gitlab_runner_tags        = var.gitlab_runner_tags
  name                      = var.name
}

##################
# INSTANCE GROUP #
##################

module "ec2-instance-group" {
  count  = var.fleeting_service == "ec2" && (var.fleeting_os == "linux" || var.fleeting_os == "macos") ? 1 : 0
  source = "../../internal/ec2/fleeting"

  fleeting_os     = var.fleeting_os
  scale_min       = var.scale_min
  scale_max       = var.scale_max
  idle_percentage = var.idle_percentage

  ami                    = var.ami
  instance_type          = var.instance_type
  vpc_id                 = local.vpc_id
  subnet_id              = local.subnet_id
  asg_storage_type       = var.asg_storage_type
  asg_storage_size       = var.asg_storage_size
  asg_storage_throughput = var.asg_storage_throughput

  macos_cores_per_license              = var.macos_cores_per_license
  macos_required_license_count_per_asg = var.macos_required_license_count_per_asg

  labels = var.labels
  name   = var.name
}

###################
# RUNNER MANAGERS #
###################

module "ec2-managers" {
  count  = var.manager_service == "ec2" ? 1 : 0
  source = "../../internal/ec2/manager"

  runner_token = var.runner_token != "" ? var.runner_token : module.gitlab[0].runner_token
  executor     = var.executor
  gitlab_url   = var.gitlab_url

  fleeting_service                           = var.fleeting_service
  fleeting_service_account_access_key_id     = local.access_key_id
  fleeting_service_account_secret_access_key = local.secret_access_key

  ssh_key_pem      = local.ssh_key_pem
  ssh_key_pem_name = local.ssh_key_pem_name

  aws_asg_name = local.aws_asg_name

  capacity_per_instance = var.capacity_per_instance
  scale_min             = var.scale_min
  scale_max             = var.scale_max
  vpc_id                = local.vpc_id
  subnet_id             = local.subnet_id

  name = var.name
}

#######
# VPC #
#######

module "vpc" {
  count  = var.aws_vpc_id != "" ? 0 : 1
  source = "../../internal/ec2/vpc"

  aws_vpc_cidr        = var.aws_vpc_cidr
  aws_vpc_subnet_cidr = var.aws_vpc_subnet_cidr
  aws_zone            = var.aws_zone
  labels              = var.labels
  name                = var.name
}