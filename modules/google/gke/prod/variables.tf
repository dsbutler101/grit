############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for identification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(string)

    # Minimum required feature support. See https://docs.gitlab.com/ee/policy/experiment-beta-support.html
    min_support = string
  })
}

##############
# GKE CONFIG #
##############

variable "google_zone" {
  type        = string
  description = "The Google Cloud zone in which to create your cluster"
}

variable "deletion_protection" {
  type        = bool
  description = "Set deletion protection for the cluster"
  default     = true
}

variable "node_pools" {
  type = map(any)
}

##############
# VPC CONFIG #
##############

variable "vpc" {
  type = object({
    id        = string
    subnet_id = string
  })
  description = "VPC and subnet to use. If ID is not provided, GRIT will create that resource for the cluster"

  default = {
    id        = ""
    subnet_id = ""
  }
}
