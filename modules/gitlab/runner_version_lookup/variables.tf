############
# METADATA #
############

variable "metadata" {
  type = object({
    # Unique name used for identification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(any)

    # Minimum required feature support. See https://docs.gitlab.com/ee/policy/experiment-beta-support.html
    min_support = string
  })
}

#########################
# RUNNER VERSION LOOKUP #
#########################

variable "skew" {
  description = "Determines how many runner versions behind you would like to deploy. E.g. 0 means latest, 1 means one minor version behind. Max is 2."
  type        = number
  default     = null
}

variable "runner_version" {
  description = "The runner version you would like to deploy as X.Y.Z. This will tell you if that specific version is supported and the skew."
  type        = string
  default     = null
}

variable "allow_unsupported_versions" {
  description = "Allow runner versions that are not explicitly tested and supported with this version of GRIT"
  type        = bool
  default     = false
}