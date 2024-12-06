resource "gitlab_user_runner" "grit-e2e" {
  runner_type = "project_type"
  project_id  = data.gitlab_project.grit-e2e.id

  description = "TestEndToEnd"

  tag_list = ["TestEndToEnd"]
}

resource "gitlab_user_runner" "grit-e2e-powershell" {
  runner_type = "project_type"
  project_id  = data.gitlab_project.grit-e2e.id

  description = "TestEndToEndPowerShell"

  tag_list = ["TestEndToEndPowerShell"]
}
