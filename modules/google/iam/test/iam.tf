#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

###################
# IAM TEST MODULE #
###################

module "iam" {
  source = "../internal"

  name = var.metadata.name
}
