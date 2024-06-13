############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for indentification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = optional(map(any), {})
  })
}

#################
# RUNNER CONFIG #
#################

variable "service" {
  type        = string
  description = "The system which provides infrastructure for the Runner Managers"
}

variable "executor" {
  description = "The runner's executor type"
  type        = string
  default     = "shell"
}

variable "scale_min" {
  description = "The minimum number of instances to maintain"
  type        = number
  default     = 0
}

variable "scale_max" {
  description = "The maximum number of instances to maintain"
  type        = number
  default     = 10
}

variable "max_use_count" {
  description = "The maximum number of times an instance can be used before it is scheduled for removal"
  type        = number
  default     = 10
}

variable "idle_percentage" {
  description = "The number of idle instances to maintain as a percentage of the current number of busy instances"
  type        = number
  default     = 10
}

variable "capacity_per_instance" {
  description = "The number of concurrent job each instances can run"
  type        = number
  default     = 1
}

variable "security_group_ids" {
  description = "Security groups to apply to the runner manager VMs"
  type        = list(string)
}

variable "privileged" {
  description = "When using docker - whether to run docker in privileged mode"
  type        = bool
  default     = false
}

variable "region" {
  description = "Region to deploy the runner manager to"
  type        = string
  default     = "us-east-1"
}

variable "runner_version" {
  description = "The version of gitlab-runner"
  type        = string
  default     = "16.11.1-1"
}

variable "aws_plugin_version" {
  description = "The version of gitlab-runner"
  type        = string
  default     = "0.5.0"
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
  default     = "ewogICJhZ2VudCI6IHsKICAgICJtZXRyaWNzX2NvbGxlY3Rpb25faW50ZXJ2YWwiOiA2MCwKICAgICJsb2dmaWxlIjogIi9vcHQvYXdzL2FtYXpvbi1jbG91ZHdhdGNoLWFnZW50L2xvZ3MvYW1hem9uLWNsb3Vkd2F0Y2gtYWdlbnQubG9nIiwKICAgICJkZWJ1ZyI6IGZhbHNlLAogICAgInJ1bl9hc191c2VyIjogImN3YWdlbnQiCiAgfSwKICAibG9ncyI6IHsKICAgICJsb2dzX2NvbGxlY3RlZCI6IHsKICAgICAgImZpbGVzIjogewogICAgICAgICJjb2xsZWN0X2xpc3QiOiBbCiAgICAgICAgICB7CiAgICAgICAgICAgICJmaWxlX3BhdGgiOiAiL3Zhci9sb2cvc3lzbG9nIiwKICAgICAgICAgICAgImxvZ19ncm91cF9uYW1lIjogIlJ1bm5lci1NYW5hZ2VyLUxvZ3MiLAogICAgICAgICAgICAibG9nX3N0cmVhbV9uYW1lIjogIlJ1bm5lck1hbmFnZXItU3lzbG9nLVN0cmVhbSIsCiAgICAgICAgICAgICJ0aW1lc3RhbXBfZm9ybWF0IjogIiVIOiAlTTogJVMleSViJS1kIgogICAgICAgICAgfSwKCSAgewogICAgICAgICAgICAiZmlsZV9wYXRoIjogIi92YXIvbG9nL2Nsb3VkLWluaXQtb3V0cHV0LmxvZyIsCiAgICAgICAgICAgICJsb2dfZ3JvdXBfbmFtZSI6ICJSdW5uZXItTWFuYWdlci1Mb2dzIiwKICAgICAgICAgICAgImxvZ19zdHJlYW1fbmFtZSI6ICJSdW5uZXJNYW5hZ2VyLUNsb3VkaW5pdC1TdHJlYW0iLAogICAgICAgICAgICAidGltZXN0YW1wX2Zvcm1hdCI6ICIlSDogJU06ICVTJXklYiUtZCIKICAgICAgICAgIH0KCV0KICAgICAgfQogICAgfQogIH0KfQo="
}

variable "enable_metrics_export" {
  type        = bool
  description = "Enable GitLab runner metrics export to the specified endpoint"
  default     = false
}

variable "metrics_export_endpoint" {
  type        = string
  description = "GitLab runner metrics export endpoint"
  default     = "0.0.0.0:9042"
}

##########
# GITLAB #
##########

variable "gitlab" {
  description = "Outputs from the gitlab module. Or your own"
  type = object({
    runner_token = string
    url          = optional(string, "https://gitlab.com")
  })
}

#######
# VPC #
#######

variable "vpc" {
  description = "Outputs from the vpc module. Or your own"
  type = object({
    id        = string
    subnet_id = string
  })
}

############
# FLEETING #
############

variable "fleeting" {
  description = "Outputs from the fleeting module. Or your own"
  type = object({
    autoscaling_group_name = optional(string, "")
    ssh_key_pem_name       = optional(string, "")
    ssh_key_pem            = optional(string, "")
    username               = optional(string, "ubuntu")
  })

  default = {}
}

#######
# IAM #
#######

variable "iam" {
  description = "Outputs from the iam module. Or your own"
  type = object({
    fleeting_access_key_id     = optional(string, "")
    fleeting_secret_access_key = optional(string, "")
  })
  default = {}
}
