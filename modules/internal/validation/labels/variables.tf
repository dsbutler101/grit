variable "labels" {
  description = "Additional labels to add to resources"
  type        = map(string)
  default     = {}

  validation {
    condition = alltrue([
      for key, value in var.labels :
      (length(key) <= 63 &&
        length(key) >= 1 &&
        length(value) <= 63 &&
        can(regex("^[a-z][a-z0-9_-]+$", key)) &&
      can(regex("^[a-z0-9_-]*$", value)))
    ])
    error_message = "Labels must be 63 characters or less, must only contain lowercase letters, numeric characters, underscores and dashes."
  }

  validation {
    condition     = length(var.labels) <= 64
    error_message = "A resource cannot have more than 64 labels"
  }
}