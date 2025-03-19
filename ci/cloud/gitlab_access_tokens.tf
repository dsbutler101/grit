resource "gitlab_project_access_token" "gandalf_security_scanning_tool" {
  project      = data.gitlab_project.grit.id
  name         = "gandalf-security-scanning-tool"
  access_level = "guest"

  scopes = ["api"]

  rotation_configuration = {
    expiration_days    = 365
    rotate_before_days = 30
  }
}

resource "gitlab_project_access_token" "e2e_tests_terraform" {
  project      = data.gitlab_project.grit.id
  name         = "e2e-tests-terraform"
  access_level = "maintainer"

  scopes = ["api"]

  rotation_configuration = {
    expiration_days    = 365
    rotate_before_days = 30
  }
}

resource "gitlab_project_access_token" "e2e_tests_jobs" {
  project      = data.gitlab_project.grit.id
  name         = "e2e-tests-jobs"
  access_level = "maintainer"

  scopes = ["api", "create_runner", "manage_runner"]

  rotation_configuration = {
    expiration_days    = 365
    rotate_before_days = 30
  }
}
