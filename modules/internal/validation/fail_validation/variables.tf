# This is a work around until we can produce errors with checks.
#
# https://github.com/hashicorp/terraform/issues/32289
# https://github.com/opentofu/opentofu/issues/757

variable "message" {
  type = string
  validation {
    condition     = var.message == ""
    error_message = var.message
  }
}
