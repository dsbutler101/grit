#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../../../internal/validation/support"
  use_case = "any"
  use_case_support = tomap({
    "any" = "experimental"
  })
  min_support = var.metadata.min_support
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
