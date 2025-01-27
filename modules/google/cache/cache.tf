#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
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

resource "google_storage_bucket" "cache-bucket" {
  name   = local.bucket_name
  labels = var.metadata.labels

  location = var.bucket_location

  force_destroy            = true
  public_access_prevention = "enforced"

  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      age = var.cache_object_lifetime
    }
  }
}

resource "google_storage_bucket_iam_binding" "cache-bucket" {
  bucket  = google_storage_bucket.cache-bucket.name
  role    = "roles/storage.objectAdmin"
  members = [for email in var.service_account_emails : "serviceAccount:${email}"]
}

