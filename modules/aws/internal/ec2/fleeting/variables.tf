variable "fleeting_os" {
  type        = string
  description = "The operating system for the Fleeting Runners"
}

variable "ami" {
  type        = string
  description = "The ID of the VM image"
}

variable "instance_type" {
  type        = string
  description = ""
}

variable "scale_min" {
  type = number
}

variable "scale_max" {
  type = number
}

variable "idle_percentage" {
  type = number
}

variable "name" {
  type = string
}

variable "labels" {
  type = map(any)
}

variable "vpc_id" {
  type = string
}

variable "subnet_id" {
  type = string
}

variable "asg_storage_type" {
  type = string
}

variable "asg_storage_size" {
  type = number
}

variable "asg_storage_throughput" {
  type = number
}

variable "macos_required_license_count_per_asg" {
  type = number
}

variable "macos_cores_per_license" {
  type = number
}

