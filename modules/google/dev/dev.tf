######################
# DEV INSTANCE GROUP #
######################

module "gce" {
  count       = var.fleeting_service == "gce" ? 1 : 0
  source      = "./internal/gce"
  fleeting_os = var.fleeting_os
}