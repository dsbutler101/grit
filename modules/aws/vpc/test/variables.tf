############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for indentification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(any)
  })
}

##############
# VPC CONFIG #
##############

variable "zone" {
  type        = string
  description = "TODO"
  default     = "us-east-1"
}

variable "cidr" {
  type        = string
  description = "TODO"
}

variable "subnet_cidr" {
  type        = string
  description = "TODO"
}
