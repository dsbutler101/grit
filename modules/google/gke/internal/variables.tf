variable "name" {
  type = string
}

variable "labels" {
  type = map(string)
}

variable "google_region" {
  type = string
}

variable "google_zone" {
  type = string
}

variable "nodes_count" {
  type = string
}

variable "node_machine_type" {
  type = string
}

variable "vpc" {
  type = object({
    id        = string
    subnet_id = string
  })
}
