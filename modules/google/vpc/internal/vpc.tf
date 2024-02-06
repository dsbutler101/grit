resource "google_compute_network" "default" {
  name = var.name

  auto_create_subnetworks = false
}

resource "google_compute_firewall" "runner-manager-ingress-default" {
  name    = "${var.name}-ingress-default"
  network = google_compute_network.default.id

  direction = "INGRESS"
  priority  = 65535

  deny {
    protocol = "all"
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_firewall" "runner-manager-egress-default" {
  name    = "${var.name}-egress-default"
  network = google_compute_network.default.id

  direction = "EGRESS"
  priority  = 65535

  allow {
    protocol = "all"
  }

  destination_ranges = ["0.0.0.0/0"]
}

resource "google_compute_subnetwork" "subnetwork" {
  for_each = var.subnetworks

  network = google_compute_network.default.id
  region  = var.google_region

  name          = each.key
  ip_cidr_range = each.value
}
