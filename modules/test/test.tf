locals {
  kubernetes_manager = var.manager_provider == "helm" || var.manager_provider == "operator"
}

#######################
# GITLAB REGISTRATION #
#######################

module "gitlab" {
  source = "../internal/gitlab"
}

######################
# KUBERNETES CLUSTER #
######################

module "gke-cluster" {
  count  = var.runner_provider == "gke" || kubernetes_manager ? 1 : 0
  source = "../internal/kubernetes/gke"
}

module "eks-cluster" {
  count  = var.runner_provider == "eks" || kubernetes_manager ? 1 : 0
  source = "../internal/kubernetes/eks"
}

module "aks-cluster" {
  count  = var.runner_provider == "aks" || kubernetes_manager ? 1 : 0
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
  count  = var.manager_provider == "gce"
  source = "../internal/manager/gce"
}

module "ec2-managers" {
  count  = var.manager_provider == "ec2"
  source = "../internal/manager/ec2"
}

module "azure-managers" {
  count  = var.manager_provider == "azure"
  source = "../internal/manager/azure"
}

module "helm" {
  count  = var.manager_provider == "helm"
  source = "../internal/manager/helm"
}

module "operator" {
  count  = var.manager_provider == "operator"
  source = "../internal/manager/operator"
}
