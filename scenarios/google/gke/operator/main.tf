locals {
  default_labels = {
    managed = "grit"
  }

  metadata = {
    name        = var.name
    labels      = merge(local.default_labels, var.labels)
    min_support = "experimental"
  }

  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnetwork_ids[local.metadata.name]
  }
}

module "vpc" {
  source = "../../../../modules/google/vpc/test/"

  metadata = local.metadata

  google_region = var.google_region

  subnetworks = zipmap(
    [local.metadata.name],
    [var.subnet_cidr],
  )
}

module "cluster" {
  source = "../../../../modules/google/gke/test/"

  metadata = local.metadata

  google_region = var.google_region

  google_zone = var.google_zone
  nodes_count = var.node_count

  vpc = local.vpc

  depends_on = [module.vpc]
}

module "operator" {
  source = "../../../../modules/k8s/operator/test/"

  depends_on = [module.cluster]
}

module "gitlab" {
  source = "../../../../modules/gitlab/test/"

  metadata           = local.metadata
  project_id         = var.gitlab_project_id
  runner_description = var.runner_description
}

module "runner" {
  source = "../../../../modules/k8s/runner/test/"

  metadata        = local.metadata
  namespace       = module.operator.namespace
  gitlab          = module.gitlab
  config_template = var.config_template

  depends_on = [module.operator]
}
