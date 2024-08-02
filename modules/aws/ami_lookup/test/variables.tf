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
  description = "The use case for the AMI"
  type        = string
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
