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

  validation {
    condition     = var.max_instances <= 1000
    error_message = "Fleeting plugin for Google will not allow to manage more than 1000 instances at once"
  }
}

variable "concurrent" {
  type = number

  default = 20

  validation {
    condition     = var.concurrent <= 2000
    error_message = "Configuration (especially network size) will not be able to handle more than 2000 concurrent jobs"
  }
}

variable "idle_percentage" {
  type = number

  default = 80
}

variable "autoscaling_policy" {
  type = object({
    scale_min          = optional(number, 1)
    idle_time          = optional(string, "2m0s")
    scale_factor       = optional(number, 0)
    scale_factor_limit = optional(number, 0)
  })

  default = {
    scale_min    = 1
    scale_factor = 0
  }
}

variable "ephemeral_runner" {
  type = object({
    disk_type    = optional(string, "gp3")
    disk_size    = optional(number, 25)
    machine_type = optional(string, "t3.medium")
    source_image = optional(string, "ami-0735db9b38fcbdb39")
  })

  default = {
    disk_type    = "gp3"
    disk_size    = 25
    machine_type = "t3.medium"
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
