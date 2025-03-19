resource "gitlab_project_variable" "google_credentials_b64" {
  project       = data.gitlab_project.grit.id
  key           = "GOOGLE_CREDENTIALS_B64"
  value         = google_service_account_key.grit_ci_key.private_key
  masked        = true
  raw           = true
  variable_type = "file"
  description   = "Terraform managed (ci/cloud/gitlab_variables.tf)"
}

resource "gitlab_project_variable" "google_project" {
  project     = data.gitlab_project.grit.id
  key         = "GOOGLE_PROJECT"
  value       = local.google_project
  masked      = true
  raw         = true
  description = "Terraform managed (ci/cloud/gitlab_variables.tf)"
}

resource "gitlab_project_variable" "google_region" {
  project     = data.gitlab_project.grit.id
  key         = "GOOGLE_REGION"
  value       = local.google_region
  raw         = true
  description = "Terraform managed (ci/cloud/gitlab_variables.tf)"
}

resource "gitlab_project_variable" "google_zone" {
  project     = data.gitlab_project.grit.id
  key         = "GOOGLE_ZONE"
  value       = local.google_zone
  raw         = true
  description = "Terraform managed (ci/cloud/gitlab_variables.tf)"
}

resource "gitlab_project_variable" "aws_access_key_id" {
  project     = data.gitlab_project.grit.id
  key         = "AWS_ACCESS_KEY_ID"
  value       = aws_iam_access_key.grit_tester.id
  raw         = true
  description = "Access key for grit-tester IAM user in shared runner sandbox. Terraform managed (ci/cloud/gitlab_variables.tf)"
}

resource "gitlab_project_variable" "aws_secret_access_key" {
  project     = data.gitlab_project.grit.id
  key         = "AWS_SECRET_ACCESS_KEY"
  value       = aws_iam_access_key.grit_tester.secret
  masked      = true
  raw         = true
  description = "Secret access key for grit-tester IAM user in shared runner sandbox. Terraform managed (ci/cloud/gitlab_variables.tf)"
}

resource "gitlab_project_variable" "aws_region" {
  project     = data.gitlab_project.grit.id
  key         = "AWS_REGION"
  value       = local.aws_region
  raw         = true
  description = "Region in which to provision resources. Terraform managed (ci/cloud/gitlab_variables.tf)"
}

resource "gitlab_project_variable" "gandalf_gitlab_token" {
  project     = data.gitlab_project.grit.id
  key         = "GANDALF_GITLAB_TOKEN"
  value       = gitlab_project_access_token.gandalf_security_scanning_tool.token
  masked      = true
  raw         = true
  description = "Project access token for gitlab-org/ci-cd/runner-tools/grit with API scope, used by Gandalf InfraSec tool to comment on MRs. Terraform managed (ci/cloud/gitlab_variables.tf). Expires on ${gitlab_project_access_token.gandalf_security_scanning_tool.expires_at}"
}

resource "gitlab_project_variable" "gitlab_token_terraform" {
  project     = data.gitlab_project.grit.id
  key         = "GITLAB_TOKEN_TERRAFORM"
  value       = gitlab_project_access_token.e2e_tests_terraform.token
  masked      = true
  raw         = true
  description = "Project access token for gitlab-org/ci-cd/runner-tools/grit with API scope, used to store terraform state. Terraform managed (ci/cloud/gitlab_variables.tf). Expires on ${gitlab_project_access_token.e2e_tests_terraform.expires_at}"
}

resource "gitlab_project_variable" "gitlab_token_tests" {
  project     = data.gitlab_project.grit.id
  key         = "GITLAB_TOKEN"
  value       = gitlab_project_access_token.e2e_tests_jobs.token
  masked      = true
  raw         = true
  description = "Project access token for gitlab-org/ci-cd/runner-tools/grit with API, create_runner and manage_runner scope, used to run integration/e2e tests and manage runners they need. Terraform managed (ci/cloud/gitlab_variables.tf). Expires on ${gitlab_project_access_token.e2e_tests_jobs.expires_at}"
}
