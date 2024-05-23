variable "name" {
  type = string
}

variable "labels" {
  type = map(string)

  validation {
    # Response from the google API:
    # Error: error creating NodePool: googleapi: Error 400: Invalid field
    # 'cluster.node_config.labels.value': "hhoerl@gitlab.com". It must begin and end
    # with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_),
    # dots (.), and alphanumerics between.
    condition = alltrue([
      for v in values(var.labels) : can(regex("^\\w+[\\w\\-\\_\\.]*\\w+$", v))
    ])
    error_message = "Label values must begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_) dots (.), and alphanumerics between."
  }

  validation {
    # Response from the google API:
    # Error: googleapi: Error 400: Invalid label key for hhoerl@gitlab.com:
    # name part must consist of alphanumeric characters, '-', '_' or '.', and
    # must start and end with an alphanumeric character (e.g. 'MyName',  or
    # 'my.name',  or '123-abc', regex used for validation is
    # '([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]').
    condition = alltrue([
      for k in keys(var.labels) : can(regex("([A-Za-z0-9][-A-Za-z0-9_\\.]*)?[A-Za-z0-9]", k))
    ])
    error_message = "Label keys must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character"
  }
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

variable "deletion_protection" {
  type = bool
}
