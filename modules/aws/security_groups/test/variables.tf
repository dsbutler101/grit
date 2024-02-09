############
# METADATA #
############

variable "metadata" {
  type = object({

    # Unique name used for indentification and partitioning resources
    name = string

    # Labels to apply to supported resources
    labels = optional(map(any), {})
  })
}

variable "vpc_id" {
  type        = string
  description = "The ID of the VPC"
}
