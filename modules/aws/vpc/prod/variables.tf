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
  type        = string
  description = "TODO"
}

variable "cidr" {
  type        = string
  description = "TODO"
}

variable "subnet_cidr" {
  type        = string
  description = "TODO"
}
