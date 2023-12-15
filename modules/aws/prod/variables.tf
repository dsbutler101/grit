##########
# GITLAB #
##########

variable "gitlab_project_id" {
  type        = string
  description = "The project ID in which to register the runner"
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

variable "runner_token" {
  type    = string
  default = ""
}

#################
# RUNNER CONFIG #
#################

variable "manager_service" {
  type        = string
  description = "The system which provides infrastructure for the Runner Managers"
}

variable "fleeting_service" {
  type        = string
  description = "The system which provides infrastructure for the Runners"
}

variable "fleeting_os" {
  type = string
}

variable "executor" {
  type    = string
  default = "docker-autoscaler"
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

variable "capacity_per_instance" {
  type = number
}

#############
# AWS INFRA #
#############

variable "ami" {
  type = string
}

variable "instance_type" {
  type = string
}

variable "asg_storage_type" {
  type    = string
  default = "gp3"
}

variable "asg_storage_size" {
  type    = number
  default = 500
}

variable "asg_storage_throughput" {
  type    = number
  default = 750 #must be in range of (125 - 1000)
}

variable "aws_zone" {
  type    = string
  default = "us-east-1a"
}

variable "aws_vpc_id" {
  type    = string
  default = ""
}

variable "aws_vpc_subnet_id" {
  type    = string
  default = ""
}

variable "aws_vpc_cidr" {
  type    = string
  default = ""
}

variable "aws_vpc_subnet_cidr" {
  type    = string
  default = ""
}

variable "macos_required_license_count_per_asg" {
  type    = number
  default = 20
}

variable "macos_cores_per_license" {
  type    = number
  default = 8
}

####################
# GENERAL SETTINGS #
####################

variable "name" {
  type = string

  validation {
    condition     = can(regex("^[0-9A-Za-z_-]+$", var.name))
    error_message = "For the name value only a-z, A-Z, 0-9, -, and _ are allowed."
  }

  validation {
    condition     = length(var.name) < 31
    error_message = "The name value must be 30 characters or less in length."
  }
}

variable "labels" {
  type = map(any)
  default = {
    env = "grit"
  }
}

variable "min_maturity" {
  type = string
  validation {
    condition     = var.min_maturity == "alpha" || var.min_maturity == "beta" || var.min_maturity == "stable"
    error_message = "min_maturity must be alpha, beta or stable"
  }
}
