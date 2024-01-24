resource "gitlab_user_runner" "primary" {
  description = "${var.runner_description} ${var.name}_GRIT"
  runner_type = "project_type"
  project_id  = var.project_id
  tag_list    = var.runner_tags
  untagged    = length(var.runner_tags) == 0 ? true : false
}
