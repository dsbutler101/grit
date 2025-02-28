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
