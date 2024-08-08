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

  name                       = var.metadata.name
  labels                     = var.metadata.labels
  fleeting_inbound_sg_rules  = var.fleeting_inbound_sg_rules
  fleeting_outbound_sg_rules = var.fleeting_outbound_sg_rules
  manager_inbound_sg_rules   = var.manager_inbound_sg_rules
  manager_outbound_sg_rules  = var.manager_outbound_sg_rules

  vpc_id = var.vpc_id
}
