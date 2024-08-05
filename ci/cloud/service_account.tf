resource "google_service_account" "grit_ci" {
  account_id   = "grit-ci"
  display_name = "grit-ci"
  project      = local.google_project
  description  = "Account use to execute CI in the gitlab-org/ci-cd/runner-tools/grit repository. Terraform managed (ci/cloud/service_account.tf)"
}

resource "google_service_account_key" "grit_ci_key" {
  service_account_id = google_service_account.grit_ci.name
}

resource "google_project_iam_member" "grit_ci_iam_service_account_admin" {
  project = local.google_project
  role    = "roles/iam.serviceAccountAdmin"
  member  = google_service_account.grit_ci.member
}

resource "google_project_iam_member" "grit_ci_iam_service_account_user" {
  project = local.google_project
  role    = "roles/iam.serviceAccountUser"
  member  = google_service_account.grit_ci.member
}

resource "google_project_iam_member" "grit_ci_iam_role_admin" {
  project = local.google_project
  role    = "roles/iam.roleAdmin"
  member  = google_service_account.grit_ci.member
}

resource "google_project_iam_member" "grit_ci_resourcemanager_project_iam_admin" {
  project = local.google_project
  role    = "roles/resourcemanager.projectIamAdmin"
  member  = google_service_account.grit_ci.member
}

resource "google_project_iam_member" "grit_ci_cloudkms_admin" {
  project = local.google_project
  role    = "roles/cloudkms.admin"
  member  = google_service_account.grit_ci.member
}

resource "google_project_iam_member" "grit_ci_cloudkms_crypto_key_ecrypter_decrypter" {
  project = local.google_project
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = google_service_account.grit_ci.member
}

resource "google_project_iam_member" "grit_ci_compute_admin" {
  project = local.google_project
  role    = "roles/compute.admin"
  member  = google_service_account.grit_ci.member
}
