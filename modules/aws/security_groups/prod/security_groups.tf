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
# SECURITY GROUPS PROD MODULE #
###################

module "security_groups" {
  source = "../internal"

  name   = var.metadata.name
  labels = var.metadata.labels

  vpc_id = var.vpc_id
}
