############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for identification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = map(string)

    # Minimum required feature support. See https://docs.gitlab.com/ee/policy/experiment-beta-support.html
    min_support = string
  })
}

###################
# OPERATOR CONFIG #
###################

variable "operator_version" {
  default     = "current"
  type        = string
  description = <<-EOF
    The operator version to deploy.

    For supported version see either
    https://gitlab.com/gitlab-org/gl-openshift/gitlab-runner-operator/-/releases
    or this module's output 'supported_operator_versions'.
  EOF
}

variable "override_manifests" {
  default     = ""
  type        = string
  description = <<-EOF
    The manifests for the operator deployment. If this is set, the
    `operator_version` will be ignored, and users of this module have to provide
    the path to a yaml file, containing all needed objects to apply to the cluster.
  EOF
}
