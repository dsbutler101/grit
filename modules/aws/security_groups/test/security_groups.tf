#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

###################
# SECURITY GROUPS TEST MODULE #
###################

module "security_groups" {
  source = "../internal"

  name   = var.metadata.name
  labels = var.metadata.labels
  vpc_id = var.vpc_id
}
