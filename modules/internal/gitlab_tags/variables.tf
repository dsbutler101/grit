variable "project_id" {
  type        = string
  description = "The project ID of the GitLab project"
}

variable "order_by" {
  type        = string
  description = "Return tags ordered by name, updated or created fields"
  default     = "version"
}
