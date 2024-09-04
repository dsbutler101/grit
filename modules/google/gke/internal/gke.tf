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

resource "google_container_node_pool" "node_pool" {
  for_each = var.node_pools

  name     = format("%s-%s", var.name, each.key)
  location = var.google_zone

  cluster = google_container_cluster.primary.id
  version = data.google_container_engine_versions.gke_version.release_channel_default_version[local.release_channel]

  node_count = each.value.node_count

  node_config {
    labels       = merge(var.labels, each.value.node_config.labels)
    tags         = concat(local.node_tags, each.value.node_config.tags)
    machine_type = each.value.node_config.machine_type
    image_type   = each.value.node_config.image_type
    disk_size_gb = each.value.node_config.disk_size_gb
    disk_type    = each.value.node_config.disk_type
    oauth_scopes = each.value.node_config.oauth_scopes
    metadata     = each.value.node_config.metadata
  }
}

# Needed to provide access token in the outputs
data "google_client_config" "provider" {}
