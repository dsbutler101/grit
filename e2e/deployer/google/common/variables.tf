variable "name" {
  type = string
}

variable "gitlab_project_id" {
  type = string
}

variable "runner_tag_list" {
  type = list(string)
}

variable "runner_version" {
  type = string
}

variable "google_project" {
  type = string
}

variable "google_region" {
  type = string
}

variable "google_zone" {
  type = string
}

variable "concurrent" {
  type    = number
  default = 1
}
