#####################
# AWS configuration #
#####################

variable "aws_zone" {
  type    = string
  default = "us-east-1a"
}

variable "aws_vpc_cidr" {
  type = string
}

#######################################
# AWS Autoscaling Group configuration #
#######################################

variable "required_license_count_per_asg" {
  type    = number
  default = 20
}

variable "cores_per_license" {
  type    = number
  default = 8
}

variable "asg_storage" {
  type = object({
    size       = optional(number, 500)
    type       = optional(string, "gp2")
    throughput = optional(number)
  })
}

variable "autoscaling_groups" {
  type = map(object({
    ami_id        = string
    instance_type = string
    subnet_cidr   = string
  }))

  /*
    Example usage:

    autoscaling_groups = {
      group-1 = {
        ami_id        = var.ami
        instance_type = "mac2.metal"
        subnet_cidr   = "10.0.22.0/21"
      },
      group-2 = {...},
      (...)
    }
  */
}

variable "protect_from_scale_in" {
  type    = bool
  default = true
}

variable "labels" {
  type = map(any)
  default = {
    env = "grit"
  }
}

variable "autoscaling_group_max" {
  type    = number
  default = 20
}
