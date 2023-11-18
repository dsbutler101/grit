##########################
# DEV EC2 INSTANCE GROUP #
##########################

module "macos" {
  count  = var.fleeting_os == "macos" ? 1 : 0
  source = "./macos"
  name   = var.name
}

module "linux" {
  count  = var.fleeting_os == "linux" ? 1 : 0
  source = "./linux"
  name   = var.name
}

module "windows" {
  count  = var.fleeting_os == "windows" ? 1 : 0
  source = "./windows"
  name   = var.name
}