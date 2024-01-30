variable "gitlab" {
  type = object({
    runner_token = string
    url          = string
  })
}

variable "fleeting" {
  type = object({
    autoscaling_group_name = string
    ssh_key_pem_name       = string
    ssh_key_pem            = string
  })
}

variable "iam" {
  type = object({
    fleeting_access_key_id     = string
    fleeting_secret_access_key = string
  })
}

variable "vpc" {
  description = "Outputs from the vpc module. Or your own"
  type = object({
    id        = string
    subnet_id = string
  })
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
  type = string
}

variable "labels" {
  type = map(any)
}

variable "name" {
  type = string
}
