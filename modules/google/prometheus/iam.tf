locals {
  name_underscore = replace(var.metadata.name, "-", "_")
}

resource "google_project_iam_custom_role" "prometheus_server" {
  role_id = "${local.name_underscore}_prometheusServer"
  title   = "Role for ${var.metadata.name} runner manager"

  permissions = [
    "compute.instances.list"
  ]
}

resource "google_project_iam_member" "prometheus_server" {
  project = var.google_project
  member  = "serviceAccount:${var.service_account_email}"
  role    = google_project_iam_custom_role.prometheus_server.id
}
