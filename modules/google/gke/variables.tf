############
# METADATA #
############

variable "metadata" {
  type = object({
    # Unique name used for identification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(string)

    # Minimum required feature support. See https://docs.gitlab.com/ee/policy/experiment-beta-support.html
    min_support = string
  })

  # TODO: apply this validation to all Google modules by creating and
  # using a google/internal/validation/labels module with these
  # semantics. If any should apply GRIT-wide, also create and use an
  # intneral/validation/labels module and apply to other providers.
  validation {
    # Response from the google API:
    # Error: error creating NodePool: googleapi: Error 400: Invalid field
    # 'cluster.node_config.labels.value': "hhoerl@gitlab.com". It must begin and end
    # with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_),
    # dots (.), and alphanumerics between.
    condition = alltrue([
      for v in values(var.metadata.labels) : can(regex("^\\w+[\\w\\-\\_\\.]*\\w+$", v))
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
      for k in keys(var.metadata.labels) : can(regex("([A-Za-z0-9][-A-Za-z0-9_\\.]*)?[A-Za-z0-9]", k))
    ])
    error_message = "Label keys must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character"
  }
}

module "validate_support" {
  source   = "../../internal/validation/support"
  use_case = "gke"
  use_case_support = tomap({
    "gke" = "experimental"
  })
  min_support = var.metadata.min_support
}

variable "name" {
  description = "The cluster name"
  type        = string
  default     = ""
}

##############
# GKE CONFIG #
##############

variable "google_zone" {
  type        = string
  description = "The Google Cloud zone in which to create your cluster"
}

variable "deletion_protection" {
  type        = bool
  description = "Set deletion protection for the cluster"
  default     = true
}

# TODO: should we have a top-level autoscaling module for producing
# well-known good configurations? This would become a moduled
# dependency.
variable "autoscaling" {
  type = object({
    enabled                     = bool
    auto_provisioning_locations = list(string)
    autoscaling_profile         = string
    resource_limits = list(object({
      resource_type = string
      minimum       = number
      maximum       = number
    }))
  })

  default = {
    enabled                     = false
    auto_provisioning_locations = []
    autoscaling_profile         = ""
    resource_limits             = []
  }
}

###########################
# NODE POOL CONFIGURATION #
###########################

variable "node_pools" {
  description = "The configuration required for each node pool added to the GKE cluster"
  type = map(object({
    node_count = optional(number, 0)
    autoscaling = optional(object({
      min_node_count = number
      max_node_count = number
    }))
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
      taints = optional(list(object({
        key    = optional(string)
        value  = optional(string)
        effect = optional(string)
      })))
    }), {})
  }))

  validation {
    condition     = length(var.node_pools) > 0
    error_message = "at least one node pool needs to be configured"
  }
}

module "validate_image_type" {
  source = "../../internal/validation/is_one_of"

  for_each = var.node_pools

  value   = lower(each.value.node_config.image_type)
  allowed = ["cos_containerd", "cos", "ubuntu_containerd", "ubuntu", "windows_ltsc_containerd", "windows_ltsc", "windows_sac_containerd", "windows_sac"]
  prefix  = "Image type for node pool ${each.key}"
}

variable "vpc" {
  type = object({
    enabled          = bool
    id               = optional(string)
    subnetwork_ids   = optional(map(string))
    subnetwork_cidrs = optional(map(string))
  })
  default = {
    enabled = false
  }
  description = "VPC and subnet to use. If ID is not provided, GRIT will create that resource for the cluster"
}

variable "manager_subnet_name" {
  type        = string
  description = "Name of the subnetwork where runner manager is deployed"
  default     = "runner-manager"
}
