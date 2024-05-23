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

variable "google_region" {
  type        = string
  description = "The Google Cloud region in which your cluster will reside"
}

variable "google_zone" {
  type        = string
  description = "The Google Cloud zone in which to create your cluster"
}

variable "nodes_count" {
  type        = string
  description = "The number of GKE nodes in your cluster"
  default     = 3
}

variable "node_machine_type" {
  type        = string
  description = "Machine type of cluster nodes"
  default     = "n2d-standard-2"
}

variable "deletion_protection" {
  type        = bool
  description = "Set deletion protection for the cluster"
  default     = true
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
