locals {
  runner_version = "${var.runner_version_lookup.runner_version != null ? var.runner_version_lookup.runner_version : module.default_version.runner_version}-1"
}
