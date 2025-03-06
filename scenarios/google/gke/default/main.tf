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

  gitlab_credentials = {
    url          = var.gitlab_url
    runner_token = var.gitlab_runner_token
  }
}

module "vpc" {
  source = "../../../../modules/google/vpc"

  metadata = local.metadata

  google_region = var.google_region

  subnetworks = zipmap(
    [local.metadata.name],
    [var.subnet_cidr],
  )
}

module "cluster" {
  source = "../../../../modules/google/gke/prod/"

  metadata = local.metadata

  google_zone = var.google_zone
  autoscaling = var.autoscaling
  node_pools  = var.node_pools

  deletion_protection = var.deletion_protection

  vpc = local.vpc

  depends_on = [module.vpc]
}

module "operator" {
  source = "../../../../modules/k8s/operator/prod/"

  metadata           = local.metadata
  operator_version   = var.operator_version
  override_manifests = var.override_operator_manifests

  depends_on = [module.cluster]
}

module "runner" {
  source = "../../../../modules/k8s/runner/prod/"

  metadata         = local.metadata
  namespace        = module.operator.namespace
  gitlab           = local.gitlab_credentials
  concurrent       = var.concurrent
  check_interval   = var.check_interval
  config_template  = var.config_template
  envvars          = var.envvars
  pod_spec_patches = var.pod_spec_patches
  runner_image     = var.runner_image
  helper_image     = var.helper_image
  runner_opts      = var.runner_opts
  log_level        = var.log_level
  listen_address   = var.listen_address

  depends_on = [module.operator]
}
