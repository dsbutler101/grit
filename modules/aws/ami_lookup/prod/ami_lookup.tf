#######################
# METADATA VALIDATION #
#######################

module "validate-support" {
  source   = "../../../internal/validation/support"
  use_case = local.use_case
  use_case_support = tomap({
    "linux-amd64-ephemeral" = "experimental"
    "linux-arm64-ephemeral" = "experimental"
  })
  min_support = var.metadata.min_support
}

##########################
# AMI LOOKUP PROD MODULE #
##########################

data "local_file" "manifest" {
  filename = local.manifest_file_path
}

locals {
  manifest = jsondecode(data.local_file.manifest.content)
}
