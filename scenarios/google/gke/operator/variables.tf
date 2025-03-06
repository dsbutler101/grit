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
}

variable "node_pools" {
  type = map(any)
}

variable "name" {
  description = "The name for the cluster, the runner and other created infra"
  type        = string
}

variable "gitlab_project_id" {
  description = "The GitLab project to register the runner for"
  type        = string
}

variable "runner_description" {
  description = "The description for the GitLab runner"
  type        = string
  default     = "default GitLab Runner"
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

variable "locked" {
  type        = bool
  description = "Specify whether the runner should be locked to a specific project"
  default     = false
}

variable "protected" {
  type        = bool
  description = "Specify whether the runner should only run protected branches"
  default     = false
}

variable "runner_tags" {
  type        = list(string)
  description = "List of tags to be applied to the runner"
  default     = []
}

variable "run_untagged" {
  type        = bool
  description = "Specify if jobs without tags should be run. When no runner_tags are set, it will always be true, else it will default to false"
  default     = false
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

variable "override_operator_manifests" {
  type    = string
  default = ""
}
