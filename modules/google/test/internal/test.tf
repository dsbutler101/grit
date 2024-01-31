locals {
  kubernetes_manager = var.manager_service == "helm" || var.manager_service == "operator"
}

#######################
# GITLAB REGISTRATION #
#######################

module "gitlab" {
  count                     = var.runner_token == "" && var.gitlab_project_id != "" ? 1 : 0
  source                    = "../../../gitlab/internal"
  gitlab_project_id         = var.gitlab_project_id
  gitlab_runner_description = var.gitlab_runner_description
  gitlab_runner_tags        = var.gitlab_runner_tags
  name                      = var.name
}

######################
# KUBERNETES CLUSTER #
######################

module "gke-cluster" {
  count  = var.fleeting_service == "gke" || local.kubernetes_manager ? 1 : 0
  source = "../../internal/kubernetes/fleeting/gke"
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
  count  = var.manager_service == "gce" ? 1 : 0
  source = "../../internal/gce/manager"
}

module "helm" {
  count        = var.manager_service == "helm" ? 1 : 0
  source       = "../../../helm/internal"
  runner_token = var.runner_token != "" ? var.runner_token : module.gitlab[0].runner_token
  gitlab_url   = var.gitlab_url
  name         = var.name
}

module "operator" {
  count  = var.manager_service == "operator" ? 1 : 0
  source = "../../../operator"
}
