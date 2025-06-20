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
}

############################
# K8S RUNNER CONFIGURATION #
############################

variable "gitlab" {
  description = "Outputs from the gitlab module. Or your own"
  type = object({
    runner_token = string
    url          = string
  })

  sensitive = true
}

variable "name" {
  description = "The name of the Runner. Will be used in the names of the resources created"
  type        = string
  default     = ""
}

variable "namespace" {
  description = "The namespace in which to deploy the Runner"
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
  nullable    = false
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
  type = list(object({
    name      = optional(string)
    patch     = optional(string)
    patchType = optional(string, "strategic")
  }))
  description = "A JSON or YAML format string that describes the changes which must be applied to the final PodSpec object before it is generated."

  validation {
    condition     = length([for patch in var.pod_spec_patches : true if patch.name != null && patch.name != ""]) == length(var.pod_spec_patches)
    error_message = "All pod_spec_patches must have name parameter set with a non empty value."
  }

  validation {
    condition     = length([for patch in var.pod_spec_patches : true if patch.patch != null && patch.patch != ""]) == length(var.pod_spec_patches)
    error_message = "All pod_spec_patches must have patch parameter set with a non empty value."
  }

  default = []
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
