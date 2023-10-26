######################
# DEV INSTANCE GROUP #
######################

module "ec2" {
  count       = var.fleeting_provider == "ec2" ? 1 : 0
  source      = "./ec2"
  fleeting_os = var.fleeting_os
}

module "gce" {
  count       = var.fleeting_provider == "gce" ? 1 : 0
  source      = "./gce"
  fleeting_os = var.fleeting_os
}

module "azure" {
  count       = var.fleeting_provider == "azure" ? 1 : 0
  source      = "./azure"
  fleeting_os = var.fleeting_os
}