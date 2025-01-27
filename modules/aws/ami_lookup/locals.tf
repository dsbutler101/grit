locals {
  default_manifest_file_path = "${abspath(path.module)}/manifest.json"
  manifest_file_path         = var.manifest_file != "" ? var.manifest_file : local.default_manifest_file_path
  deprecated_use_case_map    = { "aws-linux-ephemeral" = "linux-amd64-ephemeral" }
  use_case                   = var.use_case != "" ? local.deprecated_use_case_map[var.use_case] : "${var.os}-${var.arch}-${var.role}"
}
