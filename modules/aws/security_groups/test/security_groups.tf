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

  name                       = var.metadata.name
  labels                     = var.metadata.labels
  vpc_id                     = var.vpc_id
  fleeting_inbound_sg_rules  = var.fleeting_inbound_sg_rules
  fleeting_outbound_sg_rules = var.fleeting_outbound_sg_rules
  manager_inbound_sg_rules   = var.manager_inbound_sg_rules
  manager_outbound_sg_rules  = var.manager_outbound_sg_rules
}
