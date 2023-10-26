locals {
  kubernetes_manager = var.manager_provider == "helm" || var.manager_provider == "operator"
}

#######################
# GITLAB REGISTRATION #
#######################

module "gitlab" {
  source                    = "../internal/gitlab"
  gitlab_project_id         = var.gitlab_project_id
  gitlab_runner_description = var.gitlab_runner_description
  gitlab_runner_tags        = var.gitlab_runner_tags
}

######################
# KUBERNETES CLUSTER #
######################

module "gke-cluster" {
  count  = var.fleeting_provider == "gke" || local.kubernetes_manager ? 1 : 0
  source = "../internal/kubernetes/gke"
}

module "eks-cluster" {
  count  = var.fleeting_provider == "eks" || local.kubernetes_manager ? 1 : 0
  source = "../internal/kubernetes/eks"
}

module "aks-cluster" {
  count  = var.fleeting_provider == "aks" || local.kubernetes_manager ? 1 : 0
  source = "../internal/kubernetes/aks"
}

##################
# INSTANCE GROUP #
##################

module "gce-instance-group" {
  count  = var.fleeting_provider == "gce" ? 1 : 0
  source = "../internal/fleeting/gce"
}

module "ec2-instance-group" {
  count  = var.fleeting_provider == "ec2" ? 1 : 0
  source = "../internal/fleeting/ec2"

  os            = "linux"
  vm_img_id     = "ami-0a1cc31585e72ab51"
  instance_type = "t2.large"
  aws_vpc_cidr  = "10.0.0.0/24"
  scale_max     = 100
}

module "azure-instance-group" {
  count  = var.fleeting_provider == "azure" ? 1 : 0
  source = "../internal/fleeting/azure"
}

###################
# RUNNER MANAGERS #
###################

module "gce-managers" {
  count  = var.manager_provider == "gce" ? 1 : 0
  source = "../internal/manager/gce"
}

module "ec2-managers" {
  count                                      = var.manager_provider == "ec2" ? 1 : 0
  source                                     = "../internal/manager/ec2"
  runner_token                               = module.gitlab.runner_token
  executor                                   = "docker-autoscaler"
  gitlab_url                                 = var.gitlab_url
  fleeting_service_account_access_key_id     = try(module.ec2-instance-group[0].fleeting_service_account_access_key_id, "")
  fleeting_service_account_secret_access_key = try(module.ec2-instance-group[0].fleeting_service_account_secret_access_key, "")
  ssh_key_pem                                = try(module.ec2-instance-group[0].ssh_key_pem, "")
  ssh_key_pem_name                           = try(module.ec2-instance-group[0].ssh_key_pem_name, "")
  aws_asg_name                               = try(module.ec2-instance-group[0].autoscaling_group_names[0], "")
  idle_count                                 = 10
  scale_max                                  = 100
  fleeting_provider                          = var.fleeting_provider
}

module "azure-managers" {
  count  = var.manager_provider == "azure" ? 1 : 0
  source = "../internal/manager/azure"
}

module "helm" {
  count        = var.manager_provider == "helm" ? 1 : 0
  source       = "../internal/manager/helm"
  runner_token = module.gitlab.runner_token
  gitlab_url   = var.gitlab_url
}

module "operator" {
  count  = var.manager_provider == "operator" ? 1 : 0
  source = "../internal/manager/operator"
}
