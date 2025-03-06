#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

###################
# GKE TEST MODULE #
###################

module "gke" {
  source = "../internal"

  name   = var.metadata.name
  labels = var.metadata.labels

  google_zone = var.google_zone

  node_pools = var.node_pools

  deletion_protection = var.deletion_protection

  vpc = var.vpc

  autoscaling = var.autoscaling
}
