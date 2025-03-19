locals {
  ephemeral_runner_tag = "ephemeral-runner"
}

resource "google_compute_firewall" "ephemeral_runners_ssh_access" {
  name    = "${var.name}-ephemeral-runners-ssh-access"
  network = var.vpc.id
  project = var.subnetwork_project

  direction = "INGRESS"
  priority  = 1000

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.manager_subnet_cidr]

  target_tags = [local.ephemeral_runner_tag]
}

resource "google_compute_firewall" "ephemeral_runners_cross_vm_deny" {
  name    = "${var.name}-ephemeral-runners-cross-vm-deny"
  network = var.vpc.id
  project = var.subnetwork_project

  direction = "EGRESS"
  priority  = 1000

  deny {
    protocol = "all"
  }

  destination_ranges = var.cross_vm_deny_egress_destination_ranges

  target_tags = [local.ephemeral_runner_tag]
}
