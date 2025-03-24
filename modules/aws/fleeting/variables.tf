############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for indentification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(any)

    # Minimum required feature support. See https://docs.gitlab.com/ee/policy/experiment-beta-support.html
    min_support = string
  })
}

###################
# FLEETING CONFIG #
###################

variable "service" {
  description = "The AWS service on which to run jobs"
  type        = string
}

variable "os" {
  description = "The operating system to use"
  type        = string
}

variable "instance_type" {
  description = "The instance type to use in the autoscaling group"
  type        = string
}

variable "ephemeral_runner_ami" {
  description = "The machine image to use on the instances"
  type        = string
}

variable "storage_type" {
  description = "The type of the storage"
  type        = string
  default     = "gp3"
}

variable "storage_size" {
  description = "The size of the storage in GB"
  type        = number
  default     = 500
}

variable "storage_throughput" {
  description = "The throughput of the storage"
  type        = number
  default     = 750
}

variable "macos_license_count_per_asg" {
  description = "License count per ASG (MacOS only)"
  type        = number
  default     = 20
}

variable "macos_cores_per_license" {
  description = "Cores per license (MacOS only)"
  type        = number
  default     = 8
}

variable "scale_min" {
  description = "Autoscaling group minimum number of instances"
  type        = number
}

variable "scale_max" {
  description = "Autoscaling group maximum number of instances"
  type        = number
}

variable "security_group_ids" {
  description = "Security groups to apply to the fleeting VMs"
  type        = list(string)
}

variable "instance_role_profile_name" {
  description = "Instance role profile to attach to the runner manager instances"
  type        = string
  default     = null
}

variable "install_cloudwatch_agent" {
  type        = bool
  description = "Install cloudwatch agent"
  default     = false
}

variable "cloudwatch_agent_json" {
  type        = string
  description = <<EOF
    Configs of the cloudwatch agent, json formatted and base64 decoded
    ref: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch-Agent-Configuration-File-Details.html#Saving-Agent-Configuration-File
   EOF
  default     = "ewogICJhZ2VudCI6IHsKICAgICJtZXRyaWNzX2NvbGxlY3Rpb25faW50ZXJ2YWwiOiA2MCwKICAgICJsb2dmaWxlIjogIi9vcHQvYXdzL2FtYXpvbi1jbG91ZHdhdGNoLWFnZW50L2xvZ3MvYW1hem9uLWNsb3Vkd2F0Y2gtYWdlbnQubG9nIiwKICAgICJkZWJ1ZyI6IGZhbHNlLAogICAgInJ1bl9hc191c2VyIjogImN3YWdlbnQiCiAgfSwKICAibG9ncyI6IHsKICAgICJsb2dzX2NvbGxlY3RlZCI6IHsKICAgICAgImZpbGVzIjogewogICAgICAgICJjb2xsZWN0X2xpc3QiOiBbCiAgICAgICAgICB7CiAgICAgICAgICAgICJmaWxlX3BhdGgiOiAiL3Zhci9sb2cvc3lzbG9nIiwKICAgICAgICAgICAgImxvZ19ncm91cF9uYW1lIjogIkZsZWV0aW5nLUxvZ3MiLAogICAgICAgICAgICAibG9nX3N0cmVhbV9uYW1lIjogIkZsZWV0aW5nLVN5c2xvZy1TdHJlYW0iLAogICAgICAgICAgICAidGltZXN0YW1wX2Zvcm1hdCI6ICIlSDogJU06ICVTJXklYiUtZCIKICAgICAgICAgIH0sCgkgIHsKICAgICAgICAgICAgImZpbGVfcGF0aCI6ICIvdmFyL2xvZy9jbG91ZC1pbml0LW91dHB1dC5sb2ciLAogICAgICAgICAgICAibG9nX2dyb3VwX25hbWUiOiAiRmxlZXRpbmctTG9ncyIsCiAgICAgICAgICAgICJsb2dfc3RyZWFtX25hbWUiOiAiRmxlZXRpbmctQ2xvdWRpbml0LVN0cmVhbSIsCiAgICAgICAgICAgICJ0aW1lc3RhbXBfZm9ybWF0IjogIiVIOiAlTTogJVMleSViJS1kIgogICAgICAgICAgfQoJXQogICAgICB9CiAgICB9CiAgfQp9Cg=="
}

variable "mixed_instances_policy" {
  description = "Configuration block containing settings to define launch targets for Auto Scaling groups"
  type        = any
  default     = null
}

variable "ebs_encryption" {
  description = <<-EOD
  Enable EBS encryption on the volumes. Set it to true to enable encryption.
  EOD
  type        = bool
  default     = false
}

variable "kms_key_arn" {
  description = <<-EOD
  The ARN of the AWS Key Management Service (AWS KMS) customer master key (CMK) to use when creating the encrypted volume.
  EOD
  type        = string
  default     = null
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

#######
# VPC #
#######

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
