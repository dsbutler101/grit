variable "token" {
  type = string
}

variable "url" {
  type = string
}

variable "name" {
  type = string
}

variable "namespace" {
  type = string
}

################################
# RUNNER MANAGER CONFIGURATION #
################################

variable "concurrent" {
  type        = number
  description = "Number of maximum concurrent jobs handled by Runner at once"
}

variable "check_interval" {
  type        = number
  description = "Number of seconds between subsequent requests checking if GitLab has a new job for the Runner"
}

variable "locked" {
  type        = bool
  description = "Specify whether the runner should be locked to a specific project"
}

variable "protected" {
  type        = bool
  description = "Specify whether the runner should only run protected branches"
}

variable "runner_tags" {
  type        = list(string)
  description = "List of tags to be applied to the runner"
}

variable "run_untagged" {
  type        = bool
  description = "Specify if jobs without tags should be run. When no runner_tags are set, it will always be true, else it will default to false"
}

variable "config_template" {
  type        = string
  description = "A config.toml template provided to configure the runner"
}

variable "envvars" {
  type        = map(string)
  description = "Map of environment variables to set for the runner"
}

variable "runner_image" {
  type        = string
  description = "The container image for the GitLab Runner manager"
}

variable "helper_image" {
  type        = string
  description = "The container image for the GitLab Runner helper"
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
}
