variable "name" {
  type = string

  default = "aws-docker-as"
}

variable "labels" {
  type = map(string)

  default = {}
}

variable "aws_region" {
  type = string

  default = "us-east-1"
}

variable "aws_zone" {
  type = string

  default = "us-east-1b"
}

variable "gitlab_url" {
  type = string

  default = "https://gitlab.com"
}

variable "gitlab_project_id" {
  type = string
}

variable "capacity_per_instance" {
  type = number

  default = 1
}

variable "max_instances" {
  type = number

  default = 20
}

variable "concurrent" {
  type = number

  default = 50
}

variable "autoscaling_policy" {
  type = object({
    scale_min          = optional(number, 1)
    scale_max          = optional(number, 10)
    idle_time          = optional(string, "2m0s")
    scale_factor       = optional(number, 0)
    scale_factor_limit = optional(number, 0)
  })

  default = {
    scale_min    = 10
    scale_max    = 20
    scale_factor = 10
  }
}

variable "ephemeral_runner" {
  type = object({
    disk_type    = optional(string, "")
    disk_size    = optional(number, 25)
    machine_type = optional(string, "t2.medium")
    source_image = optional(string, "ami-0735db9b38fcbdb39")
  })

  default = {
    disk_type    = ""
    disk_size    = 25
    machine_type = "t2.medium"
    source_image = "ami-0735db9b38fcbdb39"
  }
}

variable "runner_description" {
  type = string

  default = "example-grit-docker-autoscaler-runner"
}

variable "runner_tags" {
  type = list(string)

  default = ["grit-runner"]
}
