#######################
# METADATA VALIDATION #
#######################

module "validate_name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate_support" {
  source   = "../../internal/validation/support"
  use_case = "iam"
  use_case_support = tomap({
    "iam" = "experimental"
  })
  min_support = var.metadata.min_support
}

###################
# IAM PROD MODULE #
###################

resource "google_service_account" "default" {
  account_id   = var.metadata.name
  display_name = "Service account for ${var.metadata.name}"
}
