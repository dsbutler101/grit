locals {
  node_tags = ["gke-node", "grit-gke"]

  release_channel = "STABLE"
}

data "google_container_engine_versions" "gke_version" {
  location = var.google_zone
}

locals {
  windows_images     = ["windows_sac", "windows_ltsc", "windows_ltsc_containerd"]
  linux_node_pools   = { for key, value in var.node_pools : key => value if !contains(local.windows_images, lower(value.node_config.image_type)) }
  windows_node_pools = { for key, value in var.node_pools : key => value if contains(local.windows_images, lower(value.node_config.image_type)) }
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

  dynamic "cluster_autoscaling" {
    for_each = var.autoscaling != null && var.autoscaling.enabled ? [var.autoscaling] : []

    content {
      enabled                     = true
      auto_provisioning_locations = cluster_autoscaling.value.auto_provisioning_locations
      autoscaling_profile         = cluster_autoscaling.value.autoscaling_profile

      dynamic "resource_limits" {
        for_each = cluster_autoscaling.value.resource_limits != null ? cluster_autoscaling.value.resource_limits : []

        content {
          resource_type = resource_limits.value.resource_type
          minimum       = resource_limits.value.minimum
          maximum       = resource_limits.value.maximum
        }
      }
    }
  }
}

resource "google_container_node_pool" "linux_node_pool" {
  for_each = local.linux_node_pools

  name     = format("%s-%s", var.name, each.key)
  location = var.google_zone

  cluster    = google_container_cluster.primary.id
  version    = data.google_container_engine_versions.gke_version.release_channel_default_version[local.release_channel]
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

    dynamic "taint" {
      for_each = each.value.node_config.taints != null ? each.value.node_config.taints : []
      content {
        key    = taint.value.key
        value  = taint.value.value
        effect = taint.value.effect
      }
    }
  }

  dynamic "autoscaling" {
    for_each = each.value.autoscaling != null ? [each.value.autoscaling] : []
    content {
      min_node_count = autoscaling.value.min_node_count
      max_node_count = autoscaling.value.max_node_count
    }
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }
}

resource "google_container_node_pool" "windows_node_pool" {
  for_each = local.windows_node_pools

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

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  depends_on = [
    google_container_node_pool.linux_node_pool
  ]
}

# Needed to provide access token in the outputs
data "google_client_config" "provider" {}
