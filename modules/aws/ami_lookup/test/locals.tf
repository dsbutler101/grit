locals {
  default_manifest_file_path = "${abspath(path.module)}/../manifest.json"
  manifest_file_path         = var.manifest_file != "" ? var.manifest_file : local.default_manifest_file_path
}
