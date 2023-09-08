variable "gitlab_runner_token" {
  description = "The runner token from your GitLab instance"
}

variable "gitlab_url" {
  description = "The URL of your GitLab instance"
}

variable "project_id" {
  description = "project id"
}

variable "region" {
  description = "region"
}

variable "gke_num_nodes" {
  default     = 2
  description = "number of gke nodes"
}
