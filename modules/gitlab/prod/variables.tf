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

variable "project_id" {
  description = "The numeric project ID to which to register the runner."
  type        = string
}

variable "runner_description" {
  description = "The runner description shown in the UI."
  type        = string
}

variable "runner_tags" {
  description = "The list of runner tags for selecting jobs. An empty list will run untagged jobs."
  type        = list(string)
}
