variable "allowed" {
  type    = list(string)
  default = []
}

variable "value" {
  type    = string
  default = ""
}

variable "disable" {
  type    = bool
  default = false
}

variable "prefix" {
  type    = string
  default = ""
}

locals {
  found        = contains(var.allowed, var.value)
  fail_message = "${var.prefix == "" ? "" : "${var.prefix}: "}'${var.value}' not allowed (allowed: ${join(", ", var.allowed)})"
}

module "check-allowed" {
  source  = "../fail_validation"
  message = var.disable ? "" : local.found ? "" : local.fail_message
}
