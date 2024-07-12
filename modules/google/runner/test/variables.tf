############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for identification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(string)
  })
}

#############################
# Runner Manager deployment #
#############################

variable "google_project" {
  type        = string
  description = "Google Cloud project to use"
}

variable "google_zone" {
  type        = string
  description = "Google Cloud zone to use"
}

variable "machine_type" {
  type        = string
  description = "Machine type for runner manager instance"

  default = ""
}

variable "disk_type" {
  type        = string
  description = "Disk type to use by runner manager instance"

  default = "pd-standard"
}

variable "disk_size_gb" {
  type        = string
  description = "Disk size in GB to use by runner manager instance"

  default = 50
}

variable "service_account_email" {
  type        = string
  description = "Email of service account that will be attached to the runner manager instance"
}

variable "runner_version" {
  type        = string
  description = "Version of GitLab Runner"

  default = "v16.8.0"
}

################################
# Runner Manager configuration #
################################

variable "concurrent" {
  type        = number
  description = "Number of maximum concurrent jobs handled by Runner at once"

  default = 5
}

variable "check_interval" {
  type        = number
  description = "Number of seconds between subsequent requests checking if GitLab has a new job for the Runner"

  default = 3
}

variable "log_level" {
  type        = string
  description = "Logging level (one of: debug, info, warn, error)"

  default = "info"
}

variable "listen_address" {
  type        = string
  description = "Listener address for binding metrics and debug server to"

  default = ":9252"
}

########################
# Runner configuration #
########################

variable "gitlab_url" {
  type        = string
  description = "URL of GitLab instance to connect the Runner to"
}

variable "runner_token" {
  type        = string
  description = "Runner authentication token"
}

variable "request_concurrency" {
  type        = string
  description = "How many concurrent requests for checking new jobs can be made at once"

  default = 5
}

variable "executor" {
  type        = string
  description = "Runner executor to use"
}

variable "cache_gcs_bucket" {
  type        = string
  description = "GCS bucket name for remote cache storage"

  default = ""
}

variable "runners_global_section" {
  type        = string
  description = "Hook for injecting custom configuration of [[runners]] global section"

  default = ""
}

variable "runners_docker_section" {
  type        = string
  description = "Hook for injecting custom configuration of [runners.docker] section"

  default = ""
}

variable "default_docker_image" {
  type        = string
  description = "Default docker image to use in jobs that don't specify it explicitely"

  default = "ubuntu:latest"
}

##########################
# Fleeting configuration #
##########################

variable "fleeting_googlecompute_plugin_version" {
  type        = string
  description = "Version of fleeting-plugin-googlecompute to use"

  default = "v0.1.0"
}

variable "fleeting_instance_group_name" {
  type        = string
  description = "Instance group to use for autoscaling with fleeting"
}

variable "capacity_per_instance" {
  type        = number
  description = "Maximum number of concurrent jobs to be executed on a single autoscaled instance"

  default = 1
}

variable "max_instances" {
  type        = number
  description = "Maximum number of instances autoscaling should be able to clear"

  default = 20
}

variable "max_use_count" {
  type        = number
  description = "Number of maximum usages of an autoscaled instance before it's deleted"

  default = 1
}

variable "autoscaling_policies" {
  type = list(object({
    periods            = optional(list(string), ["* * * * *"])
    timezone           = optional(string, "")
    scale_min          = optional(number, 3)
    idle_time          = optional(string, "20m0s")
    scale_factor       = optional(number, 0)
    scale_factor_limit = optional(number, 0)
  }))
  description = "Configuration of autoscaling mechanism"

  default = []
}

#######
# VPC #
#######

variable "vpc" {
  type = object({
    id        = string
    subnet_id = string
  })
  description = "VPC and subnet to use fur runner manager deployment"
}
