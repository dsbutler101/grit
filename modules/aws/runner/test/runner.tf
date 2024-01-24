#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

######################
# RUNNER TEST MODULE #
######################

module "ec2" {
  count  = var.service == "ec2" ? 1 : 0
  source = "../internal/ec2"

  gitlab   = var.gitlab
  fleeting = var.fleeting
  iam      = var.iam
  vpc      = var.vpc

  executor              = var.executor
  capacity_per_instance = var.capacity_per_instance
  scale_min             = var.scale_min
  scale_max             = var.scale_max

  name = var.metadata.name
}
