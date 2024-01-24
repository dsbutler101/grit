############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for indentification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(any)
  })
}

###################
# FLEETING CONFIG #
###################

variable "service" {
  type        = string
  description = "The AWS service on which to run jobs"
}

variable "fleeting_os" {
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

