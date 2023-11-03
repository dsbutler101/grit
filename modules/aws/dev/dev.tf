######################
# DEV INSTANCE GROUP #
######################

module "ec2" {
  count       = var.fleeting_service == "ec2" ? 1 : 0
  source      = "./internal/ec2"
  fleeting_os = var.fleeting_os
}