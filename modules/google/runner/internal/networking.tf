resource "google_compute_firewall" "runner-manager-ssh-access" {
  name    = "${var.name}-runner-manager-ssh-access"
  network = var.vpc.id

  direction = "INGRESS"
  priority  = 1000

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]

  target_tags = [local.runner_manager_tag]
}
