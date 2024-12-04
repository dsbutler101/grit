variable "name" {
  type = string
}

variable "labels" {
  type = map(string)

  validation {
    # Response from the google API:
    # Error: error creating NodePool: googleapi: Error 400: Invalid field
    # 'cluster.node_config.labels.value': "hhoerl@gitlab.com". It must begin and end
    # with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_),
    # dots (.), and alphanumerics between.
    condition = alltrue([
      for v in values(var.labels) : can(regex("^\\w+[\\w\\-\\_\\.]*\\w+$", v))
    ])
    error_message = "Label values must begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_) dots (.), and alphanumerics between."
  }

  validation {
    # Response from the google API:
    # Error: googleapi: Error 400: Invalid label key for hhoerl@gitlab.com:
    # name part must consist of alphanumeric characters, '-', '_' or '.', and
    # must start and end with an alphanumeric character (e.g. 'MyName',  or
    # 'my.name',  or '123-abc', regex used for validation is
    # '([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]').
    condition = alltrue([
      for k in keys(var.labels) : can(regex("([A-Za-z0-9][-A-Za-z0-9_\\.]*)?[A-Za-z0-9]", k))
    ])
    error_message = "Label keys must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character"
  }
}

variable "google_region" {
  type = string
}

variable "google_zone" {
  type = string
}

###########################
# NODE POOL CONFIGURATION #
###########################

variable "node_pools" {
  type = map(object({
    node_count = optional(number, 1)
    node_config = optional(object({
      labels = optional(map(string))
      tags   = optional(list(string), [])
      oauth_scopes = optional(list(string), [
        "https://www.googleapis.com/auth/logging.write",
        "https://www.googleapis.com/auth/monitoring"
      ])
      metadata = optional(map(string), {
        disable-legacy-endpoints = "true"
      })
      machine_type = optional(string)
      image_type   = optional(string, "cos_containerd")
      disk_size_gb = optional(number)
      disk_type    = optional(string)
    }), {})
  }))

  validation {
    condition     = length(var.node_pools) > 0
    error_message = "at least one node pool needs to be configured"
  }
}

module "validate-image-type" {
  source = "../../../internal/validation/is_one_of"

  for_each = var.node_pools

  value   = lower(each.value.node_config.image_type)
  allowed = ["cos_containerd", "cos", "ubuntu_containerd", "ubuntu", "windows_ltsc_containerd", "windows_ltsc", "windows_sac_containerd", "windows_sac"]
  prefix  = "Image type for node pool ${each.key}"
}

variable "vpc" {
  type = object({
    id        = string
    subnet_id = string
  })
}

variable "deletion_protection" {
  type = bool
}
