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

##################
# SERVICE CHOICE #
##################

variable "fleeting_service" {
  type        = string
  description = "Google Cloud service on which to run jobs"
}

##########################
# FLEETING CONFIGURATION #
##########################

variable "google_project" {
  type        = string
  description = "Google Cloud project to use"
}

variable "google_zone" {
  type        = string
  description = "The zone that instances in this group should be created in"
}

variable "service_account_email" {
  type        = string
  description = "Email of service account that will be allowed to manage instances through the instance group"
}

variable "machine_type" {
  type        = string
  description = "Machine type to use for autoscaled ephemeral instances"
}

variable "disk_type" {
  type        = string
  description = "Disk type to use by autoscaled ephemeral instances"

  default = "pd-standard"
}

variable "disk_size_gb" {
  type        = string
  description = "Disk size in GB to use by autoscaled ephemeral instances"

  default = 25
}

variable "source_image" {
  type        = string
  description = "Image to use for ephemeral instances"

  default = "projects/cos-cloud/global/images/family/cos-stable"
}

#######
# VPC #
#######

variable "vpc" {
  type = object({
    id        = string
    subnet_id = string
  })
  description = "VPC and subnet to use"
}

variable "manager_subnet_cidr" {
  type        = string
  description = "CIDR of the subnetwork where runner manager is deployed"
}
