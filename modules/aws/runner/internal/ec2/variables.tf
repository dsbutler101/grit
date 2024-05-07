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
    id        = string
    subnet_id = string
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
