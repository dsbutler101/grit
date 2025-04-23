variable "name" {
  type = string
}

variable "labels" {
  type = map(string)
}

#
# Fleeting configuration
#

variable "google_project" {
  description = "Google Cloud project to use"
  type        = string
}

variable "subnetwork_project" {
  description = "Project where the subnetwork is located"
  type        = string
}

variable "google_zone" {
  type = string
}

variable "service_account_email" {
  type = string
}

variable "machine_type" {
  type = string
}

variable "disk_type" {
  type = string
}

variable "disk_size_gb" {
  type = string
}

variable "source_image" {
  type = string
}

#
# VPC
#

variable "vpc" {
  type = object({
    enabled          = bool
    id               = string
    subnetwork_ids   = map(string)
    subnetwork_cidrs = map(string)
  })
}

variable "manager_subnet_name" {
  type = string
}

variable "runners_subnet_name" {
  type = string
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
