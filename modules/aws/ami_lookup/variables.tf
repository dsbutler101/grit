############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for indentification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(any)

    # Minimum required feature support. See https://docs.gitlab.com/ee/policy/experiment-beta-support.html
    min_support = string
  })
}

#################
# LOOKUP CONFIG #
#################

variable "use_case" {
  description = "The use case for the AMI. DEPRECATED: use os, arch and role instead."
  type        = string
  default     = ""
}

variable "os" {
  description = "The AMI operating system"
  type        = string
  default     = ""
  validation {
    condition     = contains(["", "linux"], var.os)
    error_message = "Variable `os` must be `linux`"
  }
}

variable "arch" {
  description = "The AMI architecture"
  type        = string
  default     = ""
  validation {
    condition     = contains(["", "amd64", "arm64"], var.arch)
    error_message = "Variable `arch` must be `amd64` or `arm64"
  }
}

variable "role" {
  description = "The role of the AMI"
  type        = string
  default     = ""
  validation {
    condition     = contains(["", "ephemeral"], var.role)
    error_message = "Variable `role` must be `ephemeral`"
  }
}

variable "region" {
  description = "The AWS region"
  type        = string
}

variable "manifest_file" {
  description = "Path to the manifest JSON file"
  type        = string
  default     = ""
}
