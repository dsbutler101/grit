resource "google_compute_instance_template" "ephemeral-runner" {
  name_prefix = "${var.name}-"

  lifecycle {
    create_before_destroy = true
  }

  tags   = [local.ephemeral_runner_tag]
  labels = var.labels

  machine_type = var.machine_type

  metadata = {
    enable-oslogin      = false
    cos-update-strategy = "update_disabled"
  }

  network_interface {
    network    = var.vpc.id
    subnetwork = var.vpc.subnet_id
    access_config {
      nat_ip = ""
    }
  }

  disk {
    disk_type    = var.disk_type
    disk_size_gb = var.disk_size_gb
    source_image = var.source_image
  }
}

resource "google_compute_instance_group_manager" "ephemeral-runners" {
  name = var.name
  zone = var.google_zone

  base_instance_name = "${var.name}-ephemeral"

  version {
    instance_template = google_compute_instance_template.ephemeral-runner.id
  }

  wait_for_instances = false
}
