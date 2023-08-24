variable "ami" {
  type = string
  default = "ami-0fcd5ff1c92b00231"
}

variable "region" {
  type    = string
  default = "us-east1"
}

variable "labels" {
  type = map
  default = {
    env = "dev"
  }
}

variable "autoscaling_group_max" {
  type = number
  default = 20
}
