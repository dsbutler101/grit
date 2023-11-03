locals {
  kubernetes_manager = var.manager_provider == "helm" || var.manager_provider == "operator"
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

######################
# KUBERNETES CLUSTER #
######################

module "gke-cluster" {
  count  = var.fleeting_service == "gke" || local.kubernetes_manager ? 1 : 0
  source = "../../internal/kubernetes/fleeting/gke"
}

module "eks-cluster" {
  count  = var.fleeting_service == "eks" || local.kubernetes_manager ? 1 : 0
  source = "../../internal/kubernetes/fleeting/eks"
}

module "aks-cluster" {
  count  = var.fleeting_service == "aks" || local.kubernetes_manager ? 1 : 0
  source = "../../internal/kubernetes/fleeting/aks"
}

##################
# INSTANCE GROUP #
##################

module "gce-instance-group" {
  count  = var.fleeting_service == "gce" ? 1 : 0
  source = "../../internal/gce"
}

###################
# RUNNER MANAGERS #
###################

module "gce-managers" {
  count  = var.manager_provider == "gce" ? 1 : 0
  source = "../../internal/gce/manager"
}

module "helm" {
  count        = var.manager_provider == "helm" ? 1 : 0
  source       = "../../../helm"
  runner_token = module.gitlab.runner_token
  gitlab_url   = var.gitlab_url
}

module "operator" {
  count  = var.manager_provider == "operator" ? 1 : 0
  source = "../../../operator"
}
