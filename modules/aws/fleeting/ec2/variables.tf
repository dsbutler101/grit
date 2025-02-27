variable "vpc" {
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

variable "os" {
  type = string
}

variable "ami" {
  type = string
}

variable "instance_type" {
  type = string
}

variable "scale_min" {
  type = number
}

variable "scale_max" {
  type = number
}

variable "name" {
  type = string
}

variable "labels" {
  type = map(any)
}

variable "storage_type" {
  type = string
}

variable "storage_size" {
  type = number
}

variable "storage_throughput" {
  type = number
}

variable "macos_license_count_per_asg" {
  type = number
}

variable "macos_cores_per_license" {
  type = number
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

variable "mixed_instances_policy" {
  type = any
}

variable "ebs_encryption" {
  type = bool
}

variable "kms_key_arn" {
  type = string
}

variable "node_exporter" {
  description = "Configuration for node_exporter"
  type = object({
    enabled = bool
    port    = optional(number)
    version = optional(string)
  })
}
