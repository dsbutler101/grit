#######################
# METADATA VALIDATION #
#######################

module "validate_name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate_support" {
  source   = "../../internal/validation/support"
  use_case = "cache"
  use_case_support = tomap({
    "cache" = "experimental"
  })
  min_support = var.metadata.min_support
}

#####################
# CACHE PROD MODULE #
#####################

locals {
  bucket_name = var.bucket_name != "" ? var.bucket_name : "${var.metadata.name}-runner-cache"
}

resource "google_storage_bucket" "cache_bucket" {
  name   = local.bucket_name
  labels = var.metadata.labels

  location = var.bucket_location

  force_destroy               = var.force_destroy
  public_access_prevention    = var.public_access_prevention
  uniform_bucket_level_access = var.uniform_bucket_level_access

  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      age = var.cache_object_lifetime
    }
  }
}

resource "google_storage_bucket_iam_binding" "cache_bucket" {
  bucket  = google_storage_bucket.cache_bucket.name
  role    = "roles/storage.objectAdmin"
  members = [for email in var.service_account_emails : "serviceAccount:${email}"]
}
