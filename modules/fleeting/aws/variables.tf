###############
# Environment #
###############

variable "realm" {
  type    = string
  default = "saas"
}

variable "env_type" {
  type = string
}

variable "shard" {
  type = string
}

variable "gl_dept" {
  type    = string
  default = "eng-dev"
}

variable "gl_dept_group" {
  type    = string
  default = "eng-dev-verify-runner"
}

variable "gl_owner_email_handle" {
  type    = string
  default = "unknown"
}

#####################
# AWS configuration #
#####################

variable "aws_zone" {
  type    = string
  default = "us-east-1a"
}

#####################
# GCP configuration #
#####################

variable "gcp_region" {
  type = string
}

#######################################
# AWS Autoscaling Group configuration #
#######################################

variable "required_license_count_per_asg" {
  type    = number
  default = 20
}

variable "cores_per_license" {
  type    = number
  default = 8
}

variable "asg_storage" {
  type = object({
    size       = optional(number, 500)
    type       = optional(string, "gp2")
    throughput = optional(number)
  })
}

variable "autoscaling_groups" {
  type = map(object({
    ami_id        = optional(string, "ami-034ccb74da463ebe1")
    instance_type = optional(string, "mac2.metal")
    subnet_cidr   = string
  }))

  /*
    Example usage:

    autoscaling_groups = {
      saas-macos-m1-blue-1 = {
        ami_id        = "ami-034ccb74da463ebe1"
        instance_type = "mac2.metal"
        subnet_cidr   = "10.0.22.0/21"
      },
      saas-macos-m1-blue-2 = {...},
      (...)
    }
  */
}

variable "protect_from_scale_in" {
  type    = bool
  default = true
}

################
# Cache bucket #
################

variable "cache_bucket_name" {
  type = string
}

##############
# Networking #
##############

variable "aws_vpc_cidr" {
  type = string
}

variable "gcp_runner_manager_vpc_link" {
  type = string
}
