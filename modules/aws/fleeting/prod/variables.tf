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

variable "ami" {
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
  default     = 750 #must be in range of (125 - 1000)
}

variable "macos_required_license_count_per_asg" {
  description = "Required license count per ASG (MacOS only)"
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

