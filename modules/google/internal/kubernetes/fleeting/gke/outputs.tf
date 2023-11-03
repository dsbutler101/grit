output "output_map" {
  value = tomap({
    "name"           = google_container_cluster.primary.name,
    "host"           = google_container_cluster.primary.endpoint,
    "access_token"   = data.google_client_config.provider.access_token,
    "ca_certificate" = google_container_cluster.primary.master_auth[0].cluster_ca_certificate,
  })
}