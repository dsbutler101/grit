############################
# DEV AZURE INSTANCE GROUP #
############################

module "linux" {
  count  = var.os == "linux" ? 1 : 0
  source = "./linux"
}

module "windows" {
  count  = var.os == "windows" ? 1 : 0
  source = "./windows"
}