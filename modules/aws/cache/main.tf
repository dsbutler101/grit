#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../../internal/validation/support"
  use_case = "any"
  use_case_support = tomap({
    "any" = "experimental"
  })
  min_support = var.metadata.min_support
}
