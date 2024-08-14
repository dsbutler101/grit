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
  default     = "grit-gitlab-runner"
}

variable "node_count" {
  description = "The GKE cluster's node count"
  type        = number
  default     = 1
}

variable "subnet_cidr" {
  description = "The subnet's CIDR where the GKE cluster will be deployed on"
  type        = string
  default     = "10.0.0.0/10"
}

variable "gitlab_pat" {
  description = "The personal access token for GitLab instance, to create the runner registration token"
  type        = string
  sensitive   = true
}

variable "gitlab_project_id" {
  description = "The GitLab project to register the runner for"
  type        = string
}

variable "runner_description" {
  description = "The description of the GitLab Runner instance"
  type        = string
  default     = "GRIT deployed runner on GKE"
}

variable "config_template" {
  description = "The config.toml template to use for the runner"
  type        = string
  default     = ""
}
