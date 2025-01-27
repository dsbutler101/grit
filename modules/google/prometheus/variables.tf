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

#########################
# Prometheus deployment #
#########################

variable "google_project" {
  description = "Google Cloud project where resources are deployed"
  type        = string
}

variable "google_zone" {
  description = "Google Cloud zone where resources are deployed"
  type        = string
}

variable "machine_type" {
  description = "Google Compute Engine machine type on which Prometheus server will be deployed"
  type        = string

  default = "n2d-standard-4"
}

variable "boot_disk" {
  description = "Boot disk specification for Prometheus server instance"
  type = object({
    disk_type = optional(string, "pd-ssd")
    size_gb   = optional(number, 25)
  })

  default = {}
}

variable "data_disk" {
  description = "Data disk specification for Prometheus server instance"
  type = object({
    disk_type = optional(string, "pd-ssd")
    size_gb   = optional(number, 100)
  })

  default = {}
}

variable "prometheus_version" {
  description = "Version of Prometheus to deploy"
  type        = string

  default = "v2.55.0"
}

variable "node_exporter_version" {
  description = "Version of Node Exporter to deploy"
  type        = string

  default = "v1.8.2"
}

variable "node_exporter_port" {
  description = "Port on which Node Exporter should be listening"
  type        = number

  default = 9100
}

variable "prometheus_external_labels" {
  description = "External labels to add to each metric that is written to a remote (like Grafana Mimir)"
  type        = map(string)

  default = {}
}

variable "mimir" {
  description = "Configuration of remote write to Grafana Mimir"
  type = object({
    url    = string
    tenant = optional(string, "")
  })

  default = {
    url = "" // Empty URL disables integration with Grafana Mimir
  }
}

variable "runner_manager_nodes" {
  description = "Configuration of Runner Manager nodes to be detected and scraped for metrics"
  type = object({
    filter = string

    exporter_ports = optional(object({
      runner_manager = optional(number, 9252)
      node_exporter  = optional(number, 9100)
    }), {})

    custom_relabel_configs = optional(list(object({
      target_label  = string
      source_labels = list(string)
      regex         = optional(string, "(.*)")
      replacement   = optional(string, "$1")
      action        = optional(string, "replace")
    })), [])

    instance_labels_to_include = optional(list(string), [])
  })
}

variable "service_account_email" {
  type        = string
  description = "Email of service account that will be attached to the prometheus instance"
}

variable "vpc" {
  description = "Configuration of VPC and subnet where Prometheus instance will be deployed"
  type = object({
    id        = string
    subnet_id = string
  })
}
