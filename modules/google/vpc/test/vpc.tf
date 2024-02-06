#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

###################
# VPC TEST MODULE #
###################

module "vpc" {
  source = "../internal"

  name          = var.metadata.name
  google_region = var.google_region

  subnetworks = var.subnetworks
}
