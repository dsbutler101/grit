variable "registry" {
  type        = string
  description = "Registry from which the node exporter image should be pulled"
  default     = "quay.io"
}

variable "image_path" {
  type        = string
  description = "Path of the node exporter image in the specified registry"
  default     = "prometheus/node-exporter"
}

variable "image_tag" {
  type        = string
  description = "Tag of the node exporter image to pull"
}

variable "service_name" {
  type        = string
  description = "The name of the node exporter service"
  default     = "node-exporter"
}

variable "port" {
  type        = number
  description = "The port on which the node exporter service will run"
  default     = 9100
}
