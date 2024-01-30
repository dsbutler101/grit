############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for indentification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(any)
  })
}

######################
# GITLAB TEST CONFIG #
######################

variable "url" {
  description = "The GitLab instance URL on which to register the runner."
  type        = string
  default     = "https://gitlab.com"
}

variable "project_id" {
  description = "The numeric project ID to which to register the runner."
  type        = string
}

variable "runner_description" {
  description = "The runner description shown in the UI."
  type        = string
  default     = "GRIT Runner"
}

variable "runner_tags" {
  description = "The list of runner tags for selecting jobs. An empty list will run untagged jobs."
  type        = list(string)
  default     = []
}
