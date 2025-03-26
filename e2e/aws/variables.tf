variable "name" {
  type = string
}

variable "gitlab_project_id" {
  type = string
}

variable "runner_tag" {
  type = string
}

variable "runner_version" {
  type = string
}

variable "ami_arch" {
  type    = string
  default = "amd64"
}

variable "enable_runner_wrapper" {
  type    = bool
  default = false
}
