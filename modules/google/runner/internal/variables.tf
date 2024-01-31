variable "name" {
  type = string
}

variable "labels" {
  type = map(string)
}

#
# Runner Manager deployment
#

variable "google_project" {
  type = string
}

variable "runner_version" {
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

variable "service_account_email" {
  type = string
}

#
# Runner Manager configuration
#

variable "concurrent" {
  type = number
}

variable "check_interval" {
  type = number
}

variable "log_level" {
  type = string
}

variable "listen_address" {
  type = string
}

#
# Runner configuration
#

variable "gitlab_url" {
  type = string
}

variable "runner_token" {
  type = string
}

variable "request_concurrency" {
  type = string
}

variable "cache_gcs_bucket" {
  type = string
}

variable "runners_global_section" {
  type = string
}

variable "runners_docker_section" {
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
