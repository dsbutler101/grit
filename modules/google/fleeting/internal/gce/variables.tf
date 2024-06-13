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
  type = string
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
    id        = string
    subnet_id = string
  })
}

variable "manager_subnet_cidr" {
  type = string
}
