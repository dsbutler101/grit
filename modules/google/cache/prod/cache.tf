#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../../../internal/validation/support"
  use_case = "cache"
  use_case_support = tomap({
    "cache" = "experimental"
  })
  min_support = var.metadata.min_support
}

#####################
# CACHE PROD MODULE #
#####################

module "cache" {
  source = "../internal"

  name   = var.metadata.name
  labels = var.metadata.labels

  bucket_location        = var.bucket_location
  bucket_name            = var.bucket_name
  cache_object_lifetime  = var.cache_object_lifetime
  service_account_emails = var.service_account_emails
}
