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
  description = "The AWS zone in which to create the subnet."
  type        = string
  default     = "us-east-1"
}

variable "cidr" {
  description = "The VPC CIDR."
  type        = string
  default     = "10.0.0.0/16"
}

variable "subnet_cidr" {
  description = "The subnet CIDR."
  type        = string
  default     = "10.0.0.0/24"
}
