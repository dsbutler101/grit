variable "deletion_protection" {
  description = "Set deletion protection for the cluster"
  type        = bool
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

variable "node_pools" {
  type = map(any)
}

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
}

variable "name" {
  description = "The name for the cluster, the runner and other created infra"
  type        = string
}

variable "gitlab_runner_token" {
  description = "The GitLab Runner token to use with Runner"
  type        = string
  sensitive   = true
}

variable "gitlab_url" {
  description = "The GitLab instance URL"
  type        = string
}

################################
# RUNNER MANAGER CONFIGURATION #
################################

variable "concurrent" {
  type        = number
  description = "Number of maximum concurrent jobs handled by Runner at once"
  default     = 5
}

variable "check_interval" {
  type        = number
  description = "Number of seconds between subsequent requests checking if GitLab has a new job for the Runner"
  default     = 3
}

variable "config_template" {
  type        = string
  description = "A config.toml template provided to configure the runner"
  default     = ""
}

variable "envvars" {
  type        = map(string)
  description = "The environment variables to configure for the runner"
  default     = {}
}

variable "runner_image" {
  type        = string
  description = "The container image for the GitLab Runner manager"
  default     = ""
}

variable "helper_image" {
  type        = string
  description = "The container image for the GitLab Runner helper"
  default     = ""
}

variable "pod_spec_patches" {
  type        = any
  description = "A JSON or YAML format string that describes the changes which must be applied to the final PodSpec object before it is generated."
  default     = []
}

variable "operator_version" {
  type        = string
  description = <<-EOF
    The operator version to deploy. This should be specified in semantic version format
    (e.g. 'v1.2.3') or set to 'latest' to use the most recent release.
  EOF
  default     = "latest"
}

variable "override_operator_manifests" {
  type        = string
  description = <<-EOT
    Optional path to custom operator manifests. Supports the following formats:
      - HTTP(S) URL (e.g., "https://example.com/custom-operator.yaml")
      - Local file path with "file://" prefix (e.g., "file:///path/to/operator.yaml")
      - If empty, uses the official GitLab Runner Operator manifest
  EOT
  default     = ""

  validation {
    condition     = var.override_operator_manifests == "" || can(regex("^(https?://|file://)", var.override_operator_manifests))
    error_message = "override_manifests must be empty or start with 'http://', 'https://', or 'file://'"
  }
}

variable "runner_opts" {
  type    = map(any)
  default = {}
}

variable "log_level" {
  type        = string
  description = "The log level for the GitLab Runner manager"
  default     = "info"
}

variable "listen_address" {
  type        = string
  description = "The address to listen on for the GitLab Runner manager"
  default     = "[::]:9252"
}
