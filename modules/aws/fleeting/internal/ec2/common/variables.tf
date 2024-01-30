#######################################
# AWS Autoscaling Group configuration #
#######################################

variable "storage_size" {
  type = number
}

variable "storage_type" {
  type = string
}

variable "storage_throughput" {
  type = number
}

variable "ami_id" {
  type = string
}

variable "instance_type" {
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

