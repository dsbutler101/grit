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

