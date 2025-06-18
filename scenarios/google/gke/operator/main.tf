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
    enabled          = module.vpc.enabled
    id               = module.vpc.id
    subnetwork_ids   = module.vpc.subnetwork_ids
    subnetwork_cidrs = module.vpc.subnetwork_cidrs
  }


  # Use nonsensitive for for_each iteration while keeping tokens protected  
  runners_nonsensitive = nonsensitive({
    for name, runner in var.runners : name => {
      for key, value in runner : key => value
      if key != "runner_token"
    }
  })
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
  source = "../../../../modules/google/gke/"

  metadata = local.metadata

  google_zone = var.google_zone
  autoscaling = var.autoscaling
  node_pools  = var.node_pools

  deletion_protection = var.deletion_protection

  vpc                 = local.vpc
  manager_subnet_name = local.metadata.name

  depends_on = [module.vpc]
}

module "operator" {
  source = "../../../../modules/k8s/operator"

  metadata           = local.metadata
  operator_version   = var.operator.version
  override_manifests = var.operator.override_manifests

  depends_on = [module.cluster]
}

module "runner" {
  source = "../../../../modules/k8s/runner/"

  for_each = local.runners_nonsensitive

  metadata  = local.metadata
  namespace = module.operator.namespace
  name      = each.key

  gitlab = {
    url          = each.value.url
    runner_token = var.runners[each.key].runner_token
  }

  concurrent       = each.value.concurrent
  check_interval   = each.value.check_interval
  locked           = each.value.locked
  protected        = each.value.protected
  run_untagged     = each.value.run_untagged
  runner_tags      = each.value.runner_tags
  config_template  = each.value.config_template
  envvars          = each.value.envvars
  pod_spec_patches = each.value.pod_spec_patches
  runner_image     = each.value.runner_image
  helper_image     = each.value.helper_image
  runner_opts      = each.value.runner_opts
  log_level        = each.value.log_level
  listen_address   = each.value.listen_address

  depends_on = [module.operator]
}
