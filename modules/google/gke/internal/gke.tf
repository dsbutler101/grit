locals {
  node_tags = ["gke-node", "grit-gke"]

  release_channel = "STABLE"
}

data "google_container_engine_versions" "gke_version" {
  location = var.google_zone
}

resource "google_container_cluster" "primary" {
  name     = var.name
  location = var.google_zone

  remove_default_node_pool = true
  initial_node_count       = 1

  network    = var.vpc.id
  subnetwork = var.vpc.subnet_id

  deletion_protection = var.deletion_protection

  release_channel {
    channel = local.release_channel
  }
}

resource "google_container_node_pool" "primary_nodes" {
  name     = var.name
  location = var.google_zone

  cluster = google_container_cluster.primary.id
  version = data.google_container_engine_versions.gke_version.release_channel_default_version[local.release_channel]

  node_count = var.nodes_count

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    labels = var.labels
    tags   = local.node_tags

    machine_type = var.node_machine_type

    metadata = {
      disable-legacy-endpoints = "true"
    }
  }
}

# Needed to provide access token in the outputs
data "google_client_config" "provider" {}
