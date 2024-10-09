#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

#####################
# CACHE PROD MODULE #
#####################

module "cache" {
  source = "../internal"

  labels = var.metadata.labels
  name   = var.metadata.name

  bucket_name           = var.bucket_name
  cache_object_lifetime = var.cache_object_lifetime
}
