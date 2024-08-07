variable "google_region" {
  description = "The region to deploy the into, see `gcloud compute zones`"
  type        = string
}

variable "google_zone" {
  description = "The zone to deploy the into, see `gcloud compute zones`"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR for the subnet the GKE cluster will be deployed on"
  type        = string
  default     = "10.0.0.0/10"
}

variable "labels" {
  description = "Labels to add to the created GKE cluster"
  type        = map(string)
  default     = {}
}

variable "node_count" {
  description = "Number of nodes for the GKE cluster"
  type        = number
  default     = 1
}

variable "name" {
  description = "The name for the cluster, the runner and other created infra"
  type        = string
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
  description = "The description for the GitLab runner"
  type        = string
  default     = "default GitLab Runner"
}

variable "config_template" {
  description = "A config.toml template provided to configure the runner"
  type        = string
  default     = ""
}
