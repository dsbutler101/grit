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
