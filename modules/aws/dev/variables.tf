variable "fleeting_service" {
  type        = string
  description = "The system which providers infrastructure for the Fleeting Runners"
}

variable "fleeting_os" {
  type        = string
  description = "The operating system for the Fleeting Runners"
}

variable "name" {
  type    = string
  default = "dev_env"
}