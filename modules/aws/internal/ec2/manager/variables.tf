variable "gitlab_url" {
  type = string
}

variable "runner_token" {
  type = string
}

variable "aws_asg_name" {
  type = string
}

variable "idle_count" {
  type = string
}

variable "scale_max" {
  type = number
}

variable "executor" {
  type = string
}

variable "ssh_key_pem" {
  type    = string
  default = ""
}
variable "ssh_key_pem_name" {
  type    = string
  default = ""
}

variable "fleeting_service_account_secret_access_key" {
  type    = string
  default = ""
}

variable "fleeting_service_account_access_key_id" {
  type    = string
  default = ""
}

variable "fleeting_service" {
  type        = string
  description = "The system which provides infrastructure for the Runners"
}