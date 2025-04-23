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

variable "subnetwork_project" {
  description = "Project where the subnetwork is located"
  type        = string
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
    enabled          = bool
    id               = string
    subnetwork_ids   = map(string)
    subnetwork_cidrs = map(string)
  })
  description = "VPC and subnet to use"
}

variable "manager_subnet_name" {
  type        = string
  description = "Name of the subnetwork where runner manager is deployed"
  default     = "runner-manager"
}

variable "runners_subnet_name" {
  type        = string
  description = "Name of the subnetwork where ephemeral runners are deployed"
  default     = "ephemeral-runners"
}

variable "additional_tags" {
  type        = list(string)
  description = "Additional tags to attach to the fleeting instances"
  default     = []
}

variable "cross_vm_deny_egress_destination_ranges" {
  description = "List of destination ranges to deny egress cross-VM communication to"
  type        = list(string)
  default     = ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"]
}

variable "access_config_enabled" {
  description = "Runner manager access config enabled"
  type        = bool
  default     = true
}
