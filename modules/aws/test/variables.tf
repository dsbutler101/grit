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
  type        = string
  description = "The system which provides infrastructure for the Runners"
  default     = "linux"
}

variable "executor" {
  type    = string
  default = "docker-autoscaler"
}

variable "scale_min" {
  type    = number
  default = 0
}

variable "scale_max" {
  type    = number
  default = 10
}

variable "idle_percentage" {
  type    = number
  default = 10
}

variable "capacity_per_instance" {
  type    = number
  default = 1
}

#############
# AWS INFRA #
#############

variable "ami" {
  type    = string
  default = "ami-0a1cc31585e72ab51"
}

variable "instance_type" {
  type    = string
  default = "t2.large"
}

variable "asg_storage_type" {
  type    = string
  default = "gp2"
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
  default = "10.1.0.0/16"
}

variable "aws_vpc_subnet_cidr" {
  type    = string
  default = "10.1.0.0/24"
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
}

variable "labels" {
  type = map(any)
  default = {
    env = "grit"
  }
}
