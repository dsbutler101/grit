resource "google_service_account" "grit_ci" {
  account_id   = "grit-ci"
  display_name = "grit-ci"
  project      = local.google_project
  description  = "Account use to execute CI in the gitlab-org/ci-cd/runner-tools/grit repository. Terraform managed (ci/cloud/service_account.tf)"
}

resource "google_service_account_key" "grit_ci_key" {
  service_account_id = google_service_account.grit_ci.name
}

locals {
  grit_ci_roles = toset([
    "roles/iam.serviceAccountAdmin",
    "roles/iam.serviceAccountKeyAdmin",
    "roles/iam.serviceAccountUser",
    "roles/iam.roleAdmin",
    "roles/resourcemanager.projectIamAdmin",
    "roles/cloudkms.admin",
    "roles/cloudkms.cryptoKeyEncrypterDecrypter",
    "roles/compute.admin",
    "roles/container.admin",
  ])
}


resource "google_project_iam_member" "grit_ci_role" {
  for_each = local.grit_ci_roles

  project = local.google_project
  role    = each.value
  member  = google_service_account.grit_ci.member
}
