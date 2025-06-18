resource "google_compute_instance_template" "ephemeral_runner" {
  name_prefix = "${var.name}-"

  lifecycle {
    create_before_destroy = true
  }

  # shielded_instance_config {
  #   enable_secure_boot          = true
  #   enable_vtpm                 = true
  #   enable_integrity_monitoring = true
  # }

  tags   = concat([local.ephemeral_runner_tag], var.additional_tags)
  labels = var.labels

  machine_type = var.machine_type

  metadata = {
    enable-oslogin      = false
    cos-update-strategy = "update_disabled"
  }

  network_interface {
    network            = var.vpc.id
    subnetwork         = var.vpc.subnetwork_ids[var.runners_subnet_name]
    subnetwork_project = var.subnetwork_project

    dynamic "access_config" {
      for_each = var.access_config_enabled ? [1] : []
      content {
        nat_ip = ""
      }
    }
  }

  disk {
    disk_type    = var.disk_type
    disk_size_gb = var.disk_size_gb
    source_image = var.source_image
  }
}

resource "google_compute_instance_group_manager" "ephemeral_runners" {
  name = var.name
  zone = var.google_zone

  base_instance_name = "${var.name}-ephemeral"

  version {
    instance_template = google_compute_instance_template.ephemeral_runner.id
  }

  wait_for_instances = false
}
