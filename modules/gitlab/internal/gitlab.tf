resource "gitlab_user_runner" "primary" {
  description = "${var.gitlab_runner_description} ${var.name}_GRIT"
  runner_type = "project_type"
  project_id  = var.gitlab_project_id
  tag_list    = var.gitlab_runner_tags
  untagged    = length(var.gitlab_runner_tags) == 0 ? true : false
}
