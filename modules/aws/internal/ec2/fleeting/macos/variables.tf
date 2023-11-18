#####################
# AWS configuration #
#####################

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

variable "asg_ami_id" {
  type = string
}

variable "asg_instance_type" {
  type = string
}

variable "asg_subnet_cidr" {
  type = string
}

variable "labels" {
  type = map(any)
  default = {
    env = "grit"
  }
}

variable "name" {
  type = string
}