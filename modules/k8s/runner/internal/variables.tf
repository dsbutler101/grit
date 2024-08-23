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
  description = "Specify if jobs without tags should be run. When no runner_tags are set, it will always be true, else it will default to false"
}

variable "run_untagged" {
  type        = bool
  description = "List of comma separated tags to be applied to the runner"
}

variable "config_template" {
  type        = string
  description = "A config.toml template provided to configure the runner"
}
