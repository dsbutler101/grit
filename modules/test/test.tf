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
}

######################
# KUBERNETES CLUSTER #
######################

module "gke-cluster" {
  count  = var.runner_provider == "gke" || local.kubernetes_manager ? 1 : 0
  source = "../internal/kubernetes/gke"
}

module "eks-cluster" {
  count  = var.runner_provider == "eks" || local.kubernetes_manager ? 1 : 0
  source = "../internal/kubernetes/eks"
}

module "aks-cluster" {
  count  = var.runner_provider == "aks" || local.kubernetes_manager ? 1 : 0
  source = "../internal/kubernetes/aks"
}

##################
# INSTANCE GROUP #
##################

module "gce-instance-group" {
  count  = var.runner_provider == "gce" ? 1 : 0
  source = "../internal/fleeting/gce"
}

module "ec2-instance-group" {
  count  = var.runner_provider == "ec2" ? 1 : 0
  source = "../internal/fleeting/ec2"
}

module "azure-instance-group" {
  count  = var.runner_provider == "azure" ? 1 : 0
  source = "../internal/fleeting/azure"
}

###################
# RUNNER MANAGERS #
###################

module "gce-managers" {
  count        = var.manager_provider == "gce" ? 1 : 0
  source       = "../internal/manager/gce"
}

module "ec2-managers" {
  count  = var.manager_provider == "ec2" ? 1 : 0
  source = "../internal/manager/ec2"
  runner_token = module.gitlab.runner_token
  gitlab_url   = var.gitlab_url
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
