variable "name" {
  description = "The name for the cluster, the runner and other created infra"
  type        = string
}

variable "deletion_protection" {
  description = "Set deletion protection for the cluster"
  type        = bool
  default     = true
}

variable "google_region" {
  description = "The region to deploy the into, see `gcloud compute zones`"
  type        = string
}

variable "google_zone" {
  description = "The zone to deploy the into, see `gcloud compute zones`"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR for the subnet the GKE cluster will be deployed on"
  type        = string
  default     = "10.0.0.0/10"
}

variable "labels" {
  description = "Labels to add to the created GKE cluster"
  type        = map(string)
  default     = {}
}

# Used in examples/test-runner-gke-google
# tflint-ignore: terraform_unused_declarations
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
}

################################
# RUNNER MANAGER CONFIGURATION #
################################

variable "runners" {
  description = "GitLab Runners to deploy on the cluster"
  type = map(object({
    runner_token = string

    url             = optional(string, "https://gitlab.com")
    concurrent      = optional(number, 5)
    check_interval  = optional(number, 3)
    runner_tags     = optional(list(string), [])
    config_template = optional(string)

    locked       = optional(bool, false)
    protected    = optional(bool, false)
    run_untagged = optional(bool, false)

    envvars = optional(map(string), {})
    pod_spec_patches = optional(list(object({
      name      = string
      patch     = string
      patchType = optional(string, "strategic")
    })), [])

    runner_image = optional(string, "")
    helper_image = optional(string, "")

    runner_opts    = optional(map(any), {})
    log_level      = optional(string, "info")
    listen_address = optional(string, "[::]:9252")
  }))

  sensitive = true

  validation {
    condition     = var.runners != null && length(keys(var.runners)) > 0
    error_message = "At least one runner must be configured."
  }
}

variable "operator" {
  description = "The configuration for the operator"
  type = object({
    version            = string
    override_manifests = optional(string)
  })

  default = {
    version = "latest"
  }
}

