variable "gitlab_runner_token" {
  description = "The runner token from your GitLab instance"
}

variable "gitlab_url" {
  description = "The URL of your GitLab instance"
  default     = "https://gitlab.com"
}

variable "project_id" {
  description = "The GCP project in which to create your cluster"
}

variable "region" {
  default     = "us-central1"
  description = "The GCP region in which to create your cluster"
}

variable "gke_num_nodes" {
  description = "The number of GKE nodes in your cluster"
  default     = 2
}
