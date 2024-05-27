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

variable "labels" {
  description = "Labels to add to the created GKE cluster"
  type        = map(string)
  default     = {}
}

variable "name" {
  description = "The name for the cluster, the runner and other created infra"
  type        = string
}

variable "gitlab_pat" {
  description = "The personal access token for GitLab instance, to create the runner registration token"
  type        = string
}

variable "gitlab_project_id" {
  description = "The GitLab project to register the runner for"
  type        = string
}

