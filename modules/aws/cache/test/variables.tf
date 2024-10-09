############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for identification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(any)
  })
}

################
# CACHE CONFIG #
################

variable "cache_object_lifetime" {
  type        = number
  default     = 14
  description = "Number of days after which cache objects are automatically removed"
}

variable "bucket_name" {
  type        = string
  default     = ""
  description = "Custom name for the created bucket. If none provided one will be created from the name in metadata."
}
