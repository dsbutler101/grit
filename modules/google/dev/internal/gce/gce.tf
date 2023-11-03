##########################
# DEV GCE INSTANCE GROUP #
##########################

module "linux" {
  count  = var.fleeting_os == "linux" ? 1 : 0
  source = "linux"
}

module "windows" {
  count  = var.fleeting_os == "windows" ? 1 : 0
  source = "windows"
}