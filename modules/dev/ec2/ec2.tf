##########################
# DEV EC2 INSTANCE GROUP #
##########################

module "macos" {
  count  = var.fleeting_os == "macos" ? 1 : 0
  source = "./macos"
}

module "linux" {
  count  = var.fleeting_os == "linux" ? 1 : 0
  source = "./linux"
}

module "windows" {
  count  = var.fleeting_os == "windows" ? 1 : 0
  source = "./windows"
}