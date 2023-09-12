variable "gitlab_token" {
  description = "A token to access your GitLab instance (e.g. PAT)"
}

variable "gitlab_url" {
  description = "The URL of your GitLab instance"
  default     = "https://gitlab.com"
}

variable "gitlab_project_id" {
  description = "The project ID in which to register the runner"
}
