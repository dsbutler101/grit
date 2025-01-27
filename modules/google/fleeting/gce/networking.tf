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

resource "google_compute_firewall" "ephemeral-runners-cross-vm-deny" {
  name    = "${var.name}-ephemeral-runners-cross-vm-deny"
  network = var.vpc.id

  direction = "EGRESS"
  priority  = 1000

  deny {
    protocol = "all"
  }

  destination_ranges = [
    "10.0.0.0/8",
    "172.16.0.0/12",
    "192.168.0.0/16",
  ]

  target_tags = [local.ephemeral_runner_tag]
}
