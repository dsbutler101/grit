#######################
# METADATA VALIDATION #
#######################

module "validate_name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate_support" {
  source   = "../../internal/validation/support"
  use_case = "all"
  use_case_support = tomap({
    "all" = "experimental"
  })
  min_support = var.metadata.min_support
}

#########################
# ADDITIONAL VALIDATION #
#########################

module "check_version_or_skew" {
  source  = "../../internal/validation/fail_validation"
  message = (var.runner_version == null && var.skew != null) || (var.runner_version != null && var.skew == null) ? "" : "Either version or skew must be provided. Not both."
}

# Error when unsupported version is used and not allowed.
module "check_unsupported_version" {
  source  = "../../internal/validation/fail_validation"
  message = local.unsupported_version_error
}

#########################
# RUNNER VERSION LOOKUP #
#########################

data "local_file" "manifest" {
  filename = "${path.module}/manifest.json"
}

locals {
  manifest = jsondecode(data.local_file.manifest.content)

  # Create version_to_skew map on the fly
  version_to_skew = { for idx, version in local.manifest : version => tostring(idx) }
  # Create skew_to_version map on the fly
  skew_to_version = { for idx, version in local.manifest : tostring(idx) => version }

  # Check if the provided runner version exists in the manifest
  is_supported_version = var.runner_version == null ? true : contains(local.manifest, var.runner_version)

  # Determine skew based on the input
  # For unsupported versions return -1 as skew
  skew = tonumber(
    var.skew == null ? (
      local.is_supported_version ? local.version_to_skew[var.runner_version] : -1
    ) : var.skew
  )

  # Determine runner version based on the input
  runner_version = var.runner_version == null ? local.skew_to_version[tostring(var.skew)] : var.runner_version

  # Create a error message for unsupported versions
  unsupported_version_error = !local.is_supported_version && !var.allow_unsupported_versions ? "Runner version ${var.runner_version} is not in the supported versions list" : ""
}