variable "gitlab" {
  description = "Outputs from the gitlab module. Or your own"
  type = object({
    runner_token = string
    url          = string
  })
}

variable "metadata" {
  type = object({

    # Unique name used for identification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(string)
  })
}

variable "name_override" {
  type    = string
  default = null
}

variable "namespace" {
  type = string
}
