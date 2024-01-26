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
  description = "TODO"
  type        = string
}

variable "ami" {
  description = "TODO"
  type        = string
}

variable "instance_type" {
  description = "TODO"
  type        = string
}

variable "asg_storage_type" {
  description = "TODO"
  type        = string
  default     = "gp3"
}

variable "asg_storage_size" {
  description = "TODO"
  type        = number
  default     = 500
}

variable "asg_storage_throughput" {
  description = "TODO"
  type        = number
  default     = 750 #must be in range of (125 - 1000)
}

variable "macos_required_license_count_per_asg" {
  description = "TODO"
  type        = number
  default     = 20
}

variable "macos_cores_per_license" {
  description = "TODO"
  type        = number
  default     = 8
}

variable "scale_min" {
  description = "TODO"
  type        = number
}

variable "scale_max" {
  description = "TODO"
  type        = number
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

