resource "gitlab_user_runner" "primary" {
  description = "GRIT test GKE"
  runner_type = "project_type"
  project_id  = var.gitlab_project_id
  untagged    = true
}
