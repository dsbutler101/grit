#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

#####################
# CACHE TEST MODULE #
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
