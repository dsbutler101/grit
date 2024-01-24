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
  description = "TODO"
  type        = string
  default     = "gitlab.com"
}

variable "project_id" {
  description = "TODO"
  type        = string
}

variable "runner_description" {
  description = "TODO"
  type        = string
  default     = "GRIT"
}

variable "runner_tags" {
  description = "TODO"
  type        = list(string)
  default     = []
}
