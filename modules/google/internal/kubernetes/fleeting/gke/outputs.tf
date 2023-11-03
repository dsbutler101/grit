output "name" {
  value = google_container_cluster.primary.name
}

output "host" {
  value = google_container_cluster.primary.endpoint
}

output "access_token" {
  value = data.google_client_config.provider.access_token
}

output "ca_certificate" {
  value = google_container_cluster.primary.master_auth[0].cluster_ca_certificate
}