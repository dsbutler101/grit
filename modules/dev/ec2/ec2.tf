##########################
# DEV EC2 INSTANCE GROUP #
##########################

module "macos" {
  count  = var.os == "macos" ? 1 : 0
  source = "./macos"
}

module "linux" {
  count  = var.os == "linux" ? 1 : 0
  source = "./linux"
}

module "windows" {
  count  = var.os == "windows" ? 1 : 0
  source = "./windows"
}