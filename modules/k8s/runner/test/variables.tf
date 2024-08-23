variable "gitlab" {
  description = "Outputs from the gitlab module. Or your own"
  type = object({
    runner_token = string
    url          = string
  })
}

variable "metadata" {
  type = object({

    # Unique name used for identification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(string)
  })
}

variable "name_override" {
  type    = string
  default = null
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
  description = "Specify if jobs without tags should be run. When no runner_tags are set, it will always be true, else it will default to false"
  default     = []
}

variable "run_untagged" {
  type        = bool
  description = "List of comma separated tags to be applied to the runner"
  default     = false
}

variable "config_template" {
  type        = string
  description = "A config.toml template provided to configure the runner"
  default     = ""
}
