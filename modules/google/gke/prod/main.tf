#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../../../internal/validation/support"
  use_case = "gke"
  use_case_support = tomap({
    "gke" = "experimental"
  })
  min_support = var.metadata.min_support
}

###################
# GKE PROD MODULE #
###################

module "gke" {
  source = "../internal"

  name   = var.metadata.name
  labels = var.metadata.labels

  google_zone = var.google_zone

  node_pools = var.node_pools

  deletion_protection = var.deletion_protection

  autoscaling = var.autoscaling

  vpc = var.vpc
}
