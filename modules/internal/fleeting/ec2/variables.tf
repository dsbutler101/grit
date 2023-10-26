variable "os" {
  description = "The operating system for the Fleeting Runners"
}

variable "vm_img_id" {
  description = "The ID of the VM image"
}

variable "instance_type" {
  description = ""
}

variable "aws_vpc_cidr" {
  type    = string
  default = "10.0.0.0/24"
}

variable "gitlab_url" {
  description = "The URL of the GitLab instance where to register the Runner Manager"
  default     = "https://gitlab.com/"
}

variable "gitlab_runner_description" {
  default = "GRIT"
}

variable "gitlab_runner_tags" {
  default = []
}

variable "scale_min" {
  default = 0
}

variable "scale_max" {
  default = 1
}

variable "idle_percentage" {
  default = 10
}