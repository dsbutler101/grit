resource "gitlab_project_access_token" "e2e-tests-trigger" {
  project      = data.gitlab_project.grit-e2e.id
  name         = "e2e-tests-trigger"
  access_level = "maintainer"

  scopes = ["api"]

  rotation_configuration = {
    expiration_days    = 365
    rotate_before_days = 30
  }
}

resource "gitlab_project_access_token" "gandalf-security-scanning-tool" {
  project      = data.gitlab_project.grit.id
  name         = "gandalf-security-scanning-tool"
  access_level = "guest"

  scopes = ["api"]

  rotation_configuration = {
    expiration_days    = 365
    rotate_before_days = 30
  }
}
