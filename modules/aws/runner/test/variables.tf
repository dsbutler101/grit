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
  description = "TODO"
  type        = string
  default     = "shell"
}

variable "scale_min" {
  description = "TODO"
  type        = number
  default     = 0
}

variable "scale_max" {
  description = "TODO"
  type        = number
  default     = 10
}

variable "idle_percentage" {
  description = "TODO"
  type        = number
  default     = 10
}

variable "capacity_per_instance" {
  description = "TODO"
  type        = number
  default     = 1
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
