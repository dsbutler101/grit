#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../../../internal/validation/support"
  use_case = "vpc"
  use_case_support = tomap({
    "vpc" = "experimental"
  })
  min_support = var.metadata.min_support
}

###################
# VPC PROD MODULE #
###################

module "vpc" {
  source = "../internal"

  name          = var.metadata.name
  google_region = var.google_region

  subnetworks = var.subnetworks
}