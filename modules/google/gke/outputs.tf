output "enabled" {
  value = tobool(true)
}

output "name" {
  value       = tostring(google_container_cluster.primary.name)
  description = "Name of the created cluster"
}

output "host" {
  value       = tostring("https://${google_container_cluster.primary.endpoint}")
  description = "Host of the GKE controller"
}

output "access_token" {
  value       = tostring(data.google_client_config.provider.access_token)
  description = "Access token for the GKE controller"
  sensitive   = true
}

output "ca_certificate" {
  value       = tostring(base64decode(google_container_cluster.primary.master_auth[0].cluster_ca_certificate))
  description = "CA certificates bundle for the GKE controller"
}
