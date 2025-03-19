locals {
  name_underscore = replace(var.metadata.name, "-", "_")
}

resource "google_project_iam_custom_role" "runner_manager" {
  role_id = "${local.name_underscore}_runnerManager"
  title   = "Role for ${var.metadata.name} runner manager"

  permissions = [
    "cloudkms.cryptoKeyVersions.useToDecrypt",
    "cloudkms.cryptoKeyVersions.list",
    "iam.serviceAccounts.signBlob"
  ]
}

resource "google_project_iam_member" "runner_manager" {
  project = var.google_project
  role    = google_project_iam_custom_role.runner_manager.id
  member  = "serviceAccount:${var.service_account_email}"
}
