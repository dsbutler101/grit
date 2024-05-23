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

  google_region = var.google_region
  google_zone   = var.google_zone

  nodes_count       = var.nodes_count
  node_machine_type = var.node_machine_type

  deletion_protection = var.deletion_protection

  vpc = var.vpc
}
