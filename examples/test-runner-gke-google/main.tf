locals {
  metadata = {
    name   = var.name
    labels = var.labels
  }

  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnetwork_ids[local.metadata.name]
  }
}

module "vpc" {
  source = "../../modules/google/vpc/test/"

  metadata = local.metadata

  google_region = var.google_region

  subnetworks = zipmap(
    [local.metadata.name],
    ["10.0.0.0/10"],
  )
}

module "cluster" {
  source = "../../modules/google/gke/test/"

  metadata = local.metadata

  google_region = var.google_region

  google_zone = var.google_zone
  nodes_count = 1

  vpc = local.vpc
}

module "operator" {
  source = "../../modules/k8s/operator/test/"

  depends_on = [module.cluster]
}

module "gitlab" {
  source = "../../modules/gitlab/test/"

  metadata           = local.metadata
  project_id         = var.gitlab_project_id
  runner_description = "this is just some runner, brought to you by GRIT"
}

module "runner" {
  source = "../../modules/k8s/runner/test/"

  metadata  = local.metadata
  namespace = module.operator.namespace
  gitlab    = module.gitlab
}
