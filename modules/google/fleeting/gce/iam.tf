locals {
  name_underscore = replace(var.name, "-", "_")
}

resource "google_project_iam_custom_role" "instance_group_manager" {
  role_id = "${local.name_underscore}_instanceGroupManager"
  title   = "Role for ${var.name} instance group management"

  permissions = [
    "compute.instances.get",
    "compute.instances.setMetadata",
    "compute.instanceGroupManagers.get",
    "compute.instanceGroupManagers.list",
    "compute.instanceGroupManagers.update",
  ]
}

resource "google_project_iam_member" "instance_group_manager" {
  project = var.google_project
  role    = google_project_iam_custom_role.instance_group_manager.id
  member  = "serviceAccount:${var.service_account_email}"
}
