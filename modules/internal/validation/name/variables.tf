variable "name" {
  description = "Unique name used for identification and partitioning resources"
  type        = string

  validation {
    condition     = can(regex("^[0-9A-Za-z_-]+$", var.name))
    error_message = "For the name value only a-z, A-Z, 0-9, -, and _ are allowed."
  }

  validation {
    condition     = length(var.name) <= 20
    error_message = "The name value must be 20 characters or less in length."
  }
}
