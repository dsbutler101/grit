variable "manager_provider" {
  description = "The system which provides infrastructure for the Runner Managers"
}

variable "runner_provider" {
  description = "The system which provides infrastructure for the Runners"
}

variable "gitlab_project_id" {
  description = "The project ID in which to register the runner"
}

variable "gitlab_url" {
  description = "The URL of the GitLab instance where to register the Runner Manager"
  default     = "https://gitlab.com/"
}