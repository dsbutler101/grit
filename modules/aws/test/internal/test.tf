locals {
  access_key_id     = try(module.ec2-instance-group[0].fleeting_service_account_access_key_id, "")
  secret_access_key = try(module.ec2-instance-group[0].fleeting_service_account_secret_access_key, "")
  ssh_key_pem       = try(module.ec2-instance-group[0].ssh_key_pem, "")
  ssh_key_pem_name  = try(module.ec2-instance-group[0].ssh_key_pem_name, "")
  aws_asg_name      = try(module.ec2-instance-group[0].autoscaling_group_name, "")
}

#######################
# GITLAB REGISTRATION #
#######################

module "gitlab" {
  source                    = "../../../gitlab/internal"
  gitlab_project_id         = var.gitlab_project_id
  gitlab_runner_description = var.gitlab_runner_description
  gitlab_runner_tags        = var.gitlab_runner_tags
}

##################
# INSTANCE GROUP #
##################

module "ec2-instance-group" {
  count  = var.fleeting_service == "ec2" ? 1 : 0
  source = "../../internal/ec2/fleeting"

  os            = "linux"
  ami           = "ami-0a1cc31585e72ab51"
  instance_type = "t2.large"
  aws_vpc_cidr  = "10.0.0.0/24"
  scale_max     = 100
}

###################
# RUNNER MANAGERS #
###################

module "ec2-managers" {
  count                                      = var.manager_provider == "ec2" ? 1 : 0
  source                                     = "../../internal/ec2/manager"
  runner_token                               = module.gitlab.runner_token
  executor                                   = "docker-autoscaler"
  gitlab_url                                 = var.gitlab_url
  fleeting_service_account_access_key_id     = local.access_key_id
  fleeting_service_account_secret_access_key = local.secret_access_key
  ssh_key_pem                                = local.ssh_key_pem
  ssh_key_pem_name                           = local.ssh_key_pem_name
  aws_asg_name                               = local.aws_asg_name
  idle_count                                 = 10
  scale_max                                  = 100
  fleeting_service                           = var.fleeting_service
}