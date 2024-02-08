locals {
  ephemeral_runner_tag = "ephemeral-runner"
}

resource "google_compute_firewall" "ephemeral-runners-ssh-access" {
  name    = "${var.name}-ephemeral-runners-ssh-access"
  network = var.vpc.id

  direction = "INGRESS"
  priority  = 1000

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.manager_subnet_cidr]

  target_tags = [local.ephemeral_runner_tag]
}
