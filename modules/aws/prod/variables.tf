variable "manager_provider" {
  type        = string
  description = "The system which provides infrastructure for the Runner Managers"
}

variable "fleeting_service" {
  type        = string
  description = "The system which provides infrastructure for the Runners"
}

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
  type = list(string)
}

variable "fleeting_os" {
  type = string
}

variable "ami" {
  type = string
}

variable "instance_type" {
  type = string
}

variable "aws_vpc_cidr" {
  type = string
}

variable "capacity_per_instance" {
  type = number
}

variable "scale_min" {
  type = number
}

variable "scale_max" {
  type = number
}

variable "executor" {
  type    = string
  default = "docker-autoscaler"
}

variable "min_maturity" {
  type = string
  validation {
    condition     = var.min_maturity == "alpha" || var.min_maturity == "beta" || var.min_maturity == "stable"
    error_message = "min_maturity must be alpha, beta or stable"
  }
}