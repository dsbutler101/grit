#######################################
# AWS Autoscaling Group configuration #
#######################################

variable "storage_size" {
  type = number
}

variable "storage_type" {
  type = string
}

variable "storage_throughput" {
  type = number
}

variable "ami_id" {
  type = string
}

variable "instance_type" {
  type = string
}

variable "labels" {
  type = map(any)
}

variable "scale_min" {
  type    = number
  default = 0
}

variable "scale_max" {
  type = number
}

variable "license_arn" {
  type = string
}

variable "jobs-host-resource-group-outputs" {
  type = map(any)
}

variable "name" {
  type = string
}

variable "subnet_ids" {
  type = list(string)
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
