############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for indentification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(any)

    # Minimum required feature support. See https://docs.gitlab.com/ee/policy/experiment-beta-support.html
    min_support = string
  })
}

######################
# GITLAB PROD CONFIG #
######################

variable "url" {
  description = "The GitLab instance URL on which to register the runner."
  type        = string
}

variable "group_id" {
  description = "The numeric ID to which to register the runner for group type runners."
  type        = string
  default     = ""
}

variable "project_id" {
  description = "The numeric ID to which to register the runner for project type runners."
  type        = string
  default     = ""
}

variable "runner_description" {
  description = "The runner description shown in the UI."
  type        = string
}

variable "runner_tags" {
  description = "The list of runner tags for selecting jobs. An empty list will run untagged jobs."
  type        = list(string)
}

variable "runner_type" {
  type        = string
  description = "The scope of the runner. Valid values are: instance_type, group_type, project_type."
  default     = "project_type"
  validation {
    condition     = contains(["instance_type", "group_type", "project_type"], var.runner_type)
    error_message = "The runner_type must be one of: instance_type, group_type, project_type."
  }
}

variable "lock_project_runner" {
  type        = bool
  description = "Mark runner as 'locked=true'. Usable only with project runners"
  default     = true
}