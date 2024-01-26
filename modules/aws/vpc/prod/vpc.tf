#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../../../internal/validation/support"
  use_case = "any"
  use_case_support = tomap({
    "any" = "experimental"
  })
  min_support = var.metadata.min_support
}

###################
# VPC PROD MODULE #
###################

module "vpc" {
  source      = "../internal"
  labels      = var.metadata.labels
  name        = var.metadata.name
  zone        = var.zone
  cidr        = var.cidr
  subnet_cidr = var.subnet_cidr
}
