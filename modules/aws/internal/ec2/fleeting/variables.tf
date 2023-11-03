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

variable "aws_vpc_cidr" {
  type    = string
  default = "10.0.0.0/24"
}

variable "gitlab_url" {
  type        = string
  description = "The URL of the GitLab instance where to register the Runner Manager"
  default     = "https://gitlab.com/"
}

variable "gitlab_runner_description" {
  type    = string
  default = "GRIT"
}

variable "gitlab_runner_tags" {
  type    = list(string)
  default = []
}

variable "scale_min" {
  type    = number
  default = 0
}

variable "scale_max" {
  type    = number
  default = 1
}

variable "idle_percentage" {
  type    = number
  default = 10
}