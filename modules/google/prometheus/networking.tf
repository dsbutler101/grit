resource "google_compute_firewall" "prometheus_ssh_access" {
  name    = "${var.metadata.name}-prometheus-ssh-access"
  network = var.vpc.id

  direction = "INGRESS"
  priority  = 1000

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]

  target_tags = [local.prometheus_server_tag]
}
