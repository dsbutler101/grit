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

variable "google_zone" {
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

variable "executor" {
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
# Autoscaling configuration
#

variable "fleeting_googlecompute_plugin_version" {
  type = string
}

variable "fleeting_instance_group_name" {
  type = string
}

variable "capacity_per_instance" {
  type = number
}

variable "max_instances" {
  type = number
}

variable "max_use_count" {
  type = number
}

variable "autoscaling_policies" {
  type = list(object({
    periods            = list(string)
    timezone           = string
    scale_min          = number
    idle_time          = string
    scale_factor       = number
    scale_factor_limit = number
  }))
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
