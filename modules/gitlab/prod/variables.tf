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
  description = "TODO"
  type        = string
}

variable "project_id" {
  description = "TODO"
  type        = string
}

variable "runner_description" {
  description = "TODO"
  type        = string
}

variable "runner_tags" {
  description = "TODO"
  type        = list(string)
}
