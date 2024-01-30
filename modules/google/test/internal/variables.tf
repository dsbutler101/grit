variable "manager_service" {
  type        = string
  description = "The system which provides infrastructure for the Runner Managers"
}

variable "fleeting_service" {
  type        = string
  description = "The system which provides infrastructure for the Runners"
}

variable "gitlab_project_id" {
  type        = string
  description = "The project ID in which to register the runner"
}

variable "gitlab_url" {
  type        = string
  description = "The URL of the GitLab instance where to register the Runner Manager"
}

variable "gitlab_runner_description" {
  type = string
}

variable "gitlab_runner_tags" {
  type = list(string)
}

variable "runner_token" {
  type = string
}

variable "name" {
  type = string
}

