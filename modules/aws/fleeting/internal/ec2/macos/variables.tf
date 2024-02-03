#######################################
# AWS Autoscaling Group configuration #
#######################################

variable "license_count_per_asg" {
  type = number
}

variable "cores_per_license" {
  type = number
}

variable "labels" {
  type = map(any)
}

variable "name" {
  type = string
}
