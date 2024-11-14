variable "name" {
  type = string
}

variable "labels" {
  type = map(string)
}

#
# Google Cloud configuration
#

variable "google_project" {
  type = string
}

variable "google_zone" {
  type = string
}

#
# Instance configuration
#

variable "machine_type" {
  type = string
}

variable "boot_disk" {
  type = object({
    disk_type = string
    size_gb   = number
  })
}

variable "data_disk" {
  type = object({
    disk_type = string
    size_gb   = number
  })
}

#
# Prometheus configuration
#

variable "prometheus_version" {
  type = string
}

variable "node_exporter_version" {
  type = string
}

variable "node_exporter_port" {
  type = number
}

variable "prometheus_external_labels" {
  type = map(string)
}

variable "mimir" {
  type = object({
    url    = string
    tenant = string
  })
}

variable "runner_manager_nodes" {
  type = object({
    filter = string

    exporter_ports = object({
      runner_manager = number
      node_exporter  = number
    })

    custom_relabel_configs = optional(list(object({
      target_label  = string
      source_labels = list(string)
      regex         = optional(string)
      replacement   = optional(string)
      action        = optional(string)
    })))

    instance_labels_to_include = optional(list(string))
  })
}

variable "service_account_email" {
  type = string
}

#
# VPC
#

variable "vpc" {
  type = object({
    id        = string
    subnet_id = string
  })
}
