############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for identification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(string)

    # Minimum required feature support. See https://docs.gitlab.com/ee/policy/experiment-beta-support.html
    min_support = string
  })
}

################
# CACHE CONFIG #
################

variable "cache_object_lifetime" {
  type        = number
  description = "Number of days after which untouched cache object will be automatically removed from GCS"
  default     = 14
}

variable "bucket_location" {
  type        = string
  description = "Location of GCS bucket. It's highly recommended to keep it in sync with the region and zone"
}

variable "service_account_emails" {
  type        = list(string)
  description = "List of service account emails for which objectAdmin access to the cache bucket should be added"
  default     = []
}

variable "bucket_name" {
  type        = string
  description = "Bucket name to use. If set then automatic name derived from metadata.name is not used"
  default     = ""
}

variable "force_destroy" {
  description = "Whether to allow force destroy of the bucket"
  type        = bool
  default     = true
}

variable "public_access_prevention" {
  description = "The public access prevention configuration for this bucket"
  type        = string
  default     = "enforced"
}

variable "uniform_bucket_level_access" {
  description = "Whether to enable uniform bucket-level access"
  type        = bool
  default     = false
}
