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

##############
# VPC CONFIG #
##############

variable "zone" {
  description = "The AWS zone in which to create the subnet."
  type        = string
}

variable "cidr" {
  description = "The VPC CIDR."
  type        = string
}

variable "subnet_cidr" {
  description = "The subnet CIDR."
  type        = string
}
