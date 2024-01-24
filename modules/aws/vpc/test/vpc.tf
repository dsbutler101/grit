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
  source      = "../internal"
  labels      = var.metadata.labels
  name        = var.metadata.name
  zone        = var.zone
  cidr        = var.cidr
  subnet_cidr = var.subnet_cidr
}
