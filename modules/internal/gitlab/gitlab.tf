resource "gitlab_user_runner" "primary" {
  description = var.gitlab_runner_description
  runner_type = "project_type"
  project_id  = var.gitlab_project_id
  untagged    = true
}
