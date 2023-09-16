output "project_id" {
  value       = var.project_id
  description = "GCloud Project ID"
}

output "region" {
  value       = var.region
  description = "GCloud Region"
}

output "zone" {
  value       = var.zone
  description = "GCloud Zone"
}

output "name" {
  value       = google_container_cluster.primary.name
  description = "GKE Cluster Name"
}

output "host" {
  value       = google_container_cluster.primary.endpoint
  description = "GKE Cluster Host"
}

output "access_token" {
  value       = google_container_cluster.primary.access_token
  description = "The token used to access the GKE cluster"
}

output "ca_certificate" {
  value       = google_container_cluster.primary.master_auth[0].cluster_ca_certificate
  description = "The certificate used to verify the master's authenticity"
}
