variable "google_project" {
  description = "The google project to use"
  type        = string
}

variable "google_region" {
  description = "The region to deploy the into, see `gcloud compute zones`"
  type        = string
}

variable "google_zone" {
  description = "The zone to deploy the into, see `gcloud compute zones`"
  type        = string
}

variable "gitlab_project_id" {
  description = "The GitLab project to register the runner for"
  type        = string
}

variable "override_operator_manifests" {
  type    = string
  default = "file://./operator.k8s.yaml"
}
