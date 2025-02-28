variable "node_exporter_version" {
  type        = string
  description = "The version of the Prometheus node exporter to install"
}

variable "node_exporter_port" {
  type        = number
  description = "Port that the node exporter will listen on"
}
