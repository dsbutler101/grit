variable "vpc" {
  type = object({
    id        = string
    subnet_id = string
  })
}

variable "os" {
  type = string
}

variable "ami" {
  type = string
}

variable "instance_type" {
  type = string
}

variable "scale_min" {
  type = number
}

variable "scale_max" {
  type = number
}

variable "name" {
  type = string
}

variable "labels" {
  type = map(any)
}

variable "storage_type" {
  type = string
}

variable "storage_size" {
  type = number
}

variable "storage_throughput" {
  type = number
}

variable "macos_required_license_count_per_asg" {
  type = number
}

variable "macos_cores_per_license" {
  type = number
}

