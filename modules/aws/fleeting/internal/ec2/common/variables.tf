#######################################
# AWS Autoscaling Group configuration #
#######################################

variable "asg_storage_size" {
  type = number
}

variable "asg_storage_type" {
  type = string
}

variable "asg_storage_throughput" {
  type = number
}

variable "asg_ami_id" {
  type = string
}

variable "asg_instance_type" {
  type = string
}

variable "protect_from_scale_in" {
  type    = bool
  default = true
}

variable "labels" {
  type = map(any)
}

variable "scale_min" {
  type = number
}

variable "scale_max" {
  type = number
}

variable "license_arn" {
  type = string
}

variable "jobs-host-resource-group-outputs" {
  type = map(any)
}

variable "name" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "subnet_id" {
  type = string
}

