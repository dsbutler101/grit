variable "google_project" {
  description = "The google project to use"
  type        = string
  default     = "test-project"
}

variable "google_region" {
  description = "The region to deploy the into, see `gcloud compute zones`"
  type        = string
  default     = "us-central1"
}

variable "google_zone" {
  description = "The zone to deploy the into, see `gcloud compute zones`"
  type        = string
  default     = "us-central1-a"
}

variable "override_operator_manifests" {
  type    = string
  default = "file://./operator.k8s.yaml"
}

variable "runner_token" {
  description = "GitLab Runner registration token"
  type        = string
  sensitive   = true
  default     = "test-token"
}
