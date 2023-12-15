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
}

variable "gitlab_runner_description" {
  type = string
}

variable "gitlab_runner_tags" {
  type = list(string)
}

variable "runner_token" {
  type = string
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
  type = string
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

variable "asg_storage_type" {
  type = string
}

variable "asg_storage_size" {
  type = number
}

variable "asg_storage_throughput" {
  type = number
}

variable "instance_type" {
  type = string
}

variable "aws_zone" {
  type = string
}

variable "aws_vpc_id" {
  type = string
}

variable "aws_vpc_subnet_id" {
  type = string
}

variable "aws_vpc_cidr" {
  type = string
}

variable "aws_vpc_subnet_cidr" {
  type = string
}

variable "macos_required_license_count_per_asg" {
  type = number
}

variable "macos_cores_per_license" {
  type = number
}

####################
# GENERAL SETTINGS #
####################

variable "name" {
  type = string
}

variable "labels" {
  type = map(any)
}

