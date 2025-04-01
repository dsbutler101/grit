variable "gitlab" {
  type = object({
    runner_token = string
    url          = string
  })
}

variable "fleeting" {
  type = object({
    autoscaling_group_name = string
    ssh_key_pem_name       = string
    ssh_key_pem            = string
    username               = string
  })
}

variable "iam" {
  type = object({
    fleeting_access_key_id     = string
    fleeting_secret_access_key = string
  })
}

variable "vpc" {
  description = "Outputs from the vpc module. Or your own"
  type = object({
    id         = string
    subnet_id  = optional(string)
    subnet_ids = optional(list(string))
  })
  validation {
    condition     = (var.vpc.subnet_id != null && try(length(var.vpc.subnet_ids), 0) == 0) || (var.vpc.subnet_id == null && try(length(var.vpc.subnet_ids), 0) > 0)
    error_message = "You cannot specify both 'subnet_id' and 'subnet_ids' OR empty values for both. Only one can be provided."
  }
}

variable "s3_cache" {
  type = object({
    enabled           = bool
    server_address    = string
    bucket_name       = string
    bucket_location   = string
    access_key_id     = string
    secret_access_key = string
  })
}

variable "capacity_per_instance" {
  type = number
}

variable "scale_min" {
  type = number
}

variable "scale_max" {
  type = number
}

variable "max_use_count" {
  type = number
}

variable "idle_percentage" {
  description = "The number of idle instances to maintain as a percentage of the current number of busy instances"
  type        = number
}

variable "idle_time" {
  description = "The period of inactivity after which the runner manager will terminate an instance"
  type        = string
}

variable "executor" {
  type = string
}

variable "labels" {
  type = map(any)
}

variable "name" {
  type = string
}

variable "privileged" {
  type = bool
}

variable "region" {
  type = string
}

variable "runner_repository" {
  type = string
}

variable "runner_version" {
  type = string
}

variable "aws_plugin_version" {
  type = string
}

variable "instance_role_profile_name" {
  type = string
}

variable "security_group_ids" {
  type = list(string)
}

variable "install_cloudwatch_agent" {
  type = bool
}

variable "cloudwatch_agent_json" {
  type = string
}

variable "enable_metrics_export" {
  type        = bool
  description = "Enable GitLab runner metrics export to the specified endpoint"
}

variable "metrics_export_endpoint" {
  type        = string
  description = "GitLab runner metrics export endpoint"
}

variable "default_docker_image" {
  type = string
}

variable "runner_manager_ami" {
  description = "The machine image to use on the runner manager"
  type        = string
}

variable "usage_logger" {
  type = object({
    enabled       = optional(bool, false)
    log_dir       = optional(string)
    custom_labels = optional(map(string), {})
  })
  default = {}
}

variable "acceptable_durations" {
  type = list(object({
    periods   = list(string)
    threshold = string
    timezone  = string
  }))
  default = []
}

variable "associate_public_ip_address" {
  type = bool
}

variable "instance_type" {
  type = string
}

variable "encrypted" {
  type = bool
}

variable "kms_key_id" {
  type = string
}

variable "volume_size" {
  type = number
}

variable "volume_type" {
  type = string
}

variable "throughput" {
  type = number
}

variable "node_exporter" {
  description = "Configuration for node_exporter"
  type = object({
    enabled = optional(bool, false)
    port    = optional(number, 9100)
    version = optional(string, "0.9.0")
  })
  default = {}
}

// Context: https://gitlab.com/gitlab-org/gitlab-runner/-/issues/38216
variable "runner_wrapper" {
  description = "Configure runner wrapper for deployer automation"
  type = object({
    enabled                     = optional(bool, false)
    process_termination_timeout = optional(string, "3h")
    socket_path                 = optional(string, "unix:///var/run/gitlab-runner-wrapper.sock")
  })
  default = {}
}
