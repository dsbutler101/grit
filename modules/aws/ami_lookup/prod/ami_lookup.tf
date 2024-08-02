#######################
# METADATA VALIDATION #
#######################

module "validate-support" {
  source   = "../../../internal/validation/support"
  use_case = var.use_case
  use_case_support = tomap({
    "aws-linux-ephemeral" = "experimental"
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
