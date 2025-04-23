############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for indentification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = optional(map(any), {})

    # Minimum required feature support. See https://docs.gitlab.com/ee/policy/experiment-beta-support.html
    min_support = string
  })
}

#################
# RUNNER CONFIG #
#################

variable "service" {
  type        = string
  description = "The AWS service on which to run the runner manager"
}

variable "executor" {
  description = "The runner's executor type. See https://docs.gitlab.com/runner/executors/"
  type        = string
}

variable "runners_global_section" {
  type        = string
  description = "Hook for injecting custom configuration of [[runners]] global section"

  default = ""
}

variable "idle_time" {
  description = "The period of inactivity after which the runner manager will terminate an instance"
  type        = string
  default     = "20m0s"
}

variable "scale_min" {
  description = "The minimum number of instances to maintain"
  type        = number
  default     = -1
}

variable "scale_max" {
  description = "The maximum number of instances to maintain"
  type        = number
  default     = -1
}

variable "max_use_count" {
  description = "The maximum number of times an instance can be used before it is scheduled for removal"
  type        = number
  default     = 10
}

variable "idle_percentage" {
  description = "The number of idle instances to maintain as a percentage of the current number of busy instances"
  type        = number
  default     = -1
}

variable "capacity_per_instance" {
  description = "The number of concurrent job each instances can run"
  type        = number
  default     = 1
  validation {
    condition     = var.capacity_per_instance >= 1
    error_message = "The capacity_per_instance value must be 1 or greater"
  }
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

variable "default_docker_image" {
  type        = string
  description = "When using docker - Default image to use in jobs that don't specify it explicitly"

  default = "ubuntu:latest"
}

variable "runners_docker_section" {
  type        = string
  description = "Hook for injecting custom configuration of [runners.docker] section"

  default = ""
}

variable "region" {
  description = "Region to deploy the runner manager to"
  type        = string
  default     = "us-east-1"
}

variable "runner_repository" {
  description = "The repository of gitlab-runner packages"
  type        = string
  default     = "gitlab-runner"
}

variable "runner_version" {
  description = "The version of gitlab-runner"
  type        = string
  default     = "16.11.1-1"
}

variable "aws_plugin_version" {
  description = "The version of fleeting-plugin-aws"
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

variable "acceptable_durations" {
  type = list(object({
    periods   = optional(list(string), ["* * * * * * *"])
    threshold = string
    timezone  = optional(string, "UTC")
  }))
  default = []
}

variable "associate_public_ip_address" {
  type        = bool
  description = "Whether to associate a public IP address with an instance in a VPC."
  default     = true
}

variable "instance_type" {
  type        = string
  description = "Instance type to use for the instance."
  default     = "t2.micro"
}

variable "encrypted" {
  description = "Enable EBS encryption on the volumes. Set it to true to enable encryption."
  type        = bool
  default     = false
}

variable "kms_key_id" {
  description = "The ARN of the AWS Key Management Service (AWS KMS) customer master key (CMK) to use when creating the encrypted volume."
  type        = string
  default     = null
}

variable "volume_size" {
  description = "The size of the EBS storage in GB"
  type        = number
  default     = null
}

variable "volume_type" {
  description = "The type of the EBS storage"
  type        = string
  default     = null
}

variable "throughput" {
  description = "The throughput of the EBS storage"
  type        = number
  default     = null
}

variable "runner_manager_ami" {
  description = "The machine image to use on the runner manager"
  type        = string
  default     = ""
}

variable "node_exporter" {
  description = "Configuration for node_exporter"
  type = object({
    enabled            = bool
    write_files_config = optional(list(string), [])
    commands           = optional(list(string), [])
    port               = optional(number, 9100)
    version            = optional(string, "0.9.0")
  })
  default = {
    enabled = false
  }
}

variable "create_key_pair" {
  description = "SSH key pair to create"
  type = object({
    algorithm = optional(string, "RSA")
    size      = optional(number, 4096)
  })
  default = null
}

##########
# GITLAB #
##########

variable "gitlab" {
  description = "Outputs from the gitlab module. Or your own"
  type = object({
    enabled      = bool
    runner_token = string
    url          = string
  })
}

#######
# VPC #
#######

variable "vpc" {
  description = "Outputs from the vpc module. Or your own"
  type = object({
    enabled    = bool
    id         = string
    subnet_ids = optional(list(string))
  })
}

############
# FLEETING #
############

variable "fleeting" {
  description = "Outputs from the fleeting module. Or your own"
  type = object({
    enabled                = bool
    autoscaling_group_name = optional(string, "")
    ssh_key_pem_name       = optional(string, "")
    ssh_key_pem            = optional(string, "")
    username               = optional(string, "ubuntu")
  })

  default = {
    enabled = false
  }
}

#######
# IAM #
#######

variable "iam" {
  description = "Outputs from the iam module. Or your own"
  type = object({
    enabled                    = bool
    fleeting_access_key_id     = optional(string, "")
    fleeting_secret_access_key = optional(string, "")
  })

  default = {
    enabled = false
  }
}

#########
# CACHE #
#########

variable "cache" {
  description = "Output from the cache module. Or your own"
  type = object({
    enabled           = bool
    server_address    = optional(string, "")
    bucket_name       = optional(string, "")
    bucket_location   = optional(string, "")
    access_key_id     = optional(string, "")
    secret_access_key = optional(string, "")
  })

  default = {
    enabled = false
  }
}

################
# Usage Logger #
################
variable "usage_logger" {
  description = "Enable GitLab Runner usage logging"
  type = object({
    enabled       = bool
    log_dir       = optional(string, "/var/log/usage")
    custom_labels = optional(map(string), {})
  })
  default = {
    enabled = false
  }
}

##################
# Runner Wrapper #
##################
variable "runner_wrapper" {
  description = "Enable gitlab-runner wrapper"
  type = object({
    enabled                     = bool
    process_termination_timeout = optional(string)
    socket_path                 = optional(string)
  })
  default = {
    enabled = false
  }
}
