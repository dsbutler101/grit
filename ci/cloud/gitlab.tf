resource "gitlab_project_variable" "google_credentials_b64" {
  project       = local.gitlab_project
  key           = "GOOGLE_CREDENTIALS_B64"
  value         = google_service_account_key.grit_ci_key.private_key
  masked        = true
  variable_type = "file"
  description   = "Terraform managed (ci/cloud/gitlab.tf)"
}

resource "gitlab_project_variable" "google_project" {
  project     = local.gitlab_project
  key         = "GOOGLE_PROJECT"
  value       = local.google_project
  masked      = true
  description = "Terraform managed (ci/cloud/gitlab.tf)"
}

resource "gitlab_project_variable" "google_region" {
  project     = local.gitlab_project
  key         = "GOOGLE_REGION"
  value       = local.google_region
  description = "Terraform managed (ci/cloud/gitlab.tf)"
}

resource "gitlab_project_variable" "google_zone" {
  project     = local.gitlab_project
  key         = "GOOGLE_ZONE"
  value       = local.google_zone
  description = "Terraform managed (ci/cloud/gitlab.tf)"
}
