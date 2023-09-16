######################
# DEV INSTANCE GROUP #
######################

module "ec2" {
  count = var.fleeting_provider == "ec2" ? 1 : 0
  source = "./ec2"
  os = var.os
}

module "gce" {
  count = var.fleeting_provider == "gce" ? 1 : 0
  source = "./gce"
  os = var.os
}

module "azure" {
  count = var.fleeting_provider == "azure" ? 1 : 0
  source = "./azure"
  os = var.os
}