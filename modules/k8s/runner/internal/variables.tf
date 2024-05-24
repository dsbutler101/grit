variable "gitlab" {
  description = "Outputs from the gitlab module. Or your own"
  type = object({
    runner_token = string
    url          = string
  })
}

variable "name" {
  type = string
}

variable "namespace" {
  type = string
}
