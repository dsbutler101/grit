variable "ami" {
  type = string
}

variable "region" {
  type    = string
  default = "us-east1"
}

variable "labels" {
  type = map
  default = {}
}

variable "autoscaling_group_max" {
  type = number
  default = 20
}
