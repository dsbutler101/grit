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
# VPC CONFIG #
##############

variable "google_region" {
  type        = string
  description = "The Google Cloud region in which your cluster will reside"
}

variable "subnetworks" {
  type        = map(string)
  description = "A map of subnetwork names -> CIDRs to create in the created VPC"

  default = {}
}
