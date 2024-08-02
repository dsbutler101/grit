##########################
# AMI LOOKUP PROD MODULE #
##########################

data "local_file" "manifest" {
  filename = local.manifest_file_path
}

locals {
  manifest = jsondecode(data.local_file.manifest.content)
}
