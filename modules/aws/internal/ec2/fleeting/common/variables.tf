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

variable "asg_storage_size" {
  type    = number
  default = 500
}

variable "asg_storage_type" {
  type    = string
  default = "gp2"
}

variable "asg_storage_throughput" {
  type    = number
  default = 0
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

variable "scale_min" {
  type    = number
  default = 0
}

variable "scale_max" {
  type    = number
  default = 20
}

variable "idle_percentage" {
  type    = number
  default = 10
}

variable "license_arn" {
  type    = string
  default = ""
}

variable "jobs-host-resource-group-outputs" {
  type    = map(any)
  default = {}
}

variable "name" {
  type = string
}