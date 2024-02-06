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

variable "idle_percentage" {
  description = "The number of idle instances to maintain as a percentage of the current number of busy instances"
  type        = number
  default     = -1
}

variable "capacity_per_instance" {
  description = "The number of concurrent job each instances can run"
  type        = number
  default     = -1
}

##########
# GITLAB #
##########

variable "gitlab" {
  description = "Outputs from the gitlab module. Or your own"
  type = object({
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
    autoscaling_group_name = string
    ssh_key_pem_name       = string
    ssh_key_pem            = string
  })
}

#######
# IAM #
#######

variable "iam" {
  description = "Outputs from the iam module. Or your own"
  type = object({
    fleeting_access_key_id     = string
    fleeting_secret_access_key = string
  })
}
