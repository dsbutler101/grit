#######################
# METADATA VALIDATION #
#######################

module "validate_name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate_support" {
  source   = "../../internal/validation/support"
  use_case = "vpc"
  use_case_support = tomap({
    "vpc" = "experimental"
  })
  min_support = var.metadata.min_support
}

##################
# DEFAULT LABELS #
##################

module "labels" {
  source = "../../internal/labels"

  name              = var.metadata.name
  additional_labels = var.metadata.labels
}

###################
# VPC PROD MODULE #
###################

resource "google_compute_network" "default" {
  name = var.metadata.name

  auto_create_subnetworks = false
}

resource "google_compute_firewall" "runner_manager_ingress_default" {
  name    = "${var.metadata.name}-ingress-default"
  network = google_compute_network.default.id

  direction = "INGRESS"
  priority  = 65535

  deny {
    protocol = "all"
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_firewall" "runner_manager_egress_default" {
  name    = "${var.metadata.name}-egress-default"
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
