variable "additional_labels" {
  description = "Additional labels to add to resources"
  type        = map(string)
  default     = {}
}

variable "name" {
  description = "The name for the resource"
  type        = string
}