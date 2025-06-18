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

#############################
# Runner Manager deployment #
#############################

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
  description = "Google Cloud zone to use"
}

# TODO: make this a top-level module like AWS
variable "node_exporter" {
  description = "Configuration of Node Exporter to deploy on the runner instance"
  type = object({
    version = optional(string, "v1.8.2")
    port    = optional(number, 9100)
  })

  default = {}
}

variable "machine_type" {
  type        = string
  description = "Machine type for runner manager instance"

  default = ""
}

variable "disk_image" {
  type        = string
  description = "Disk image to use by runner manager instance"

  default = "projects/cos-cloud/global/images/family/cos-stable"
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

##################
# RUNNER VERSION #
##################

variable "runner_version_lookup" {
  description = "The version of gitlab-runner"
  type = object({
    skew           = optional(number)
    runner_version = optional(string)
  })
  default = {}
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

// DEPRECATED: use runner_metrics_listener instead
variable "listen_address" {
  type        = string
  description = "Listener address for binding metrics and debug server to (DEPRECATED: use runner_metrics_listener instead)"

  default = ""
}

variable "runner_metrics_listener" {
  type = object({
    address = optional(string, "0.0.0.0")
    port    = optional(number, 9252)
  })
  description = "TCP address and port to which runner metrics and debug server listener should be attached"

  default = {}
}

##################
# Runner Wrapper #
##################

variable "runner_wrapper" {
  type = object({
    enabled                     = optional(bool, false)
    debug                       = optional(bool, false)
    process_termination_timeout = optional(string, "3h")
  })
  description = "Enable gitlab-runner wrapper"

  default = {}
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
  type        = number
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

variable "runner_manager_additional_firewall_rules" {
  type = map(object({
    direction = string
    priority  = number
    allow = optional(list(object({
      protocol = string
      ports    = list(number)
    })), [])
    deny = optional(list(object({
      protocol = string
      ports    = list(number)
    })), [])
    source_ranges = list(string)
  }))

  default = {}
}

variable "vpc" {
  type = object({
    enabled          = bool
    id               = string
    subnetwork_ids   = map(string)
    subnetwork_cidrs = map(string)
  })
  description = "VPC and subnet to use fur runner manager deployment"
}

variable "manager_subnet_name" {
  type        = string
  description = "Name of the subnetwork where runner manager is deployed"
  default     = "runner-manager"
}

variable "source_ranges" {
  description = "Runner manager source ranges"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

# KMS
variable "kms_location" {
  description = "KMS key ring location"
  type        = string
  default     = "global"
}

variable "address_type" {
  type        = string
  description = "Type of the address to be created for the runner manager instance `INTERNAL` or `EXTERNAL`"
  default     = "EXTERNAL"
}

variable "access_config_enabled" {
  description = "Runner manager access config enabled"
  type        = bool
  default     = false
}

variable "additional_tags" {
  type        = list(string)
  description = "Additional tags to attach to the runner manager instance"
  default     = []
}

variable "runner_registry" {
  type        = string
  description = "The registry where the runner image is stored"
  default     = "registry.gitlab.com/gitlab-org/gitlab-runner"
}

variable "https_proxy" {
  description = "https proxy to use"
  type        = string
  default     = ""
}

variable "http_proxy" {
  description = "http proxy to use"
  type        = string
  default     = ""
}

variable "no_proxy" {
  description = "no proxy"
  type        = string
  default     = ""
}

variable "additional_volumes" {
  type        = list(string)
  description = "Additional volumes to mount in docker runner"
  default     = []
}
