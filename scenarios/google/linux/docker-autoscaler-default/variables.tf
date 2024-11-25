variable "name" {
  type = string
}

variable "labels" {
  type = map(string)

  default = {}
}

variable "google_project" {
  type = string
}

variable "google_region" {
  type = string
}

variable "google_zone" {
  type = string
}

variable "gitlab_url" {
  type = string
}

variable "runner_token" {
  type = string
}

variable "runner_tags" {
  type = list(string)

  default = []
}

variable "runner_machine_type" {
  type = string

  default = ""
}

variable "runner_disk_type" {
  type = string

  default = ""
}

variable "concurrent" {
  type = number

  default = 50
}

variable "runners_global_section" {
  type = string

  default = ""
}

variable "runners_docker_section" {
  type = string

  default = ""
}

variable "capacity_per_instance" {
  type = number

  default = 1
}

variable "max_instances" {
  type = number

  default = 200
}

variable "max_use_count" {
  type = number

  default = 1
}

variable "autoscaling_policies" {
  type = list(object({
    periods            = optional(list(string), ["* * * * *"])
    timezone           = optional(string, "")
    scale_min          = optional(number, 3)
    idle_time          = optional(string, "20m0s")
    scale_factor       = optional(number, 0)
    scale_factor_limit = optional(number, 0)
  }))

  default = []
}

variable "ephemeral_runner" {
  type = object({
    disk_type    = optional(string, "pd-standard")
    disk_size    = optional(number, 25)
    machine_type = optional(string, "n2d-standard-2")
    source_image = optional(string, "projects/cos-cloud/global/images/family/cos-stable")
  })

  default = {
    disk_type    = "pd-standard"
    disk_size    = 25
    machine_type = "n2d-standard-2"
    source_image = "projects/cos-cloud/global/images/family/cos-stable"
  }
}

variable "prometheus" {
  type = object({
    enabled = bool

    mimir = optional(object({
      url    = string
      tenant = optional(string, "")
    }))

    external_labels = optional(map(string))

    custom_relabel_configs = optional(list(object({
      target_label  = string
      source_labels = list(string)
      regex         = optional(string, "(.*)")
      replacement   = optional(string, "$1")
      action        = optional(string, "replace")
    })))

    instance_labels_to_include = optional(list(string), [])
  })

  default = {
    enabled = false
  }
}
