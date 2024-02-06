output "name" {
  value       = module.gke.name
  description = "Name of the created cluster"
}
output "host" {
  value       = module.gke.host
  description = "Host of the GKE controller"
}

output "access_token" {
  value       = module.gke.access_token
  description = "Access token for the GKE controller"
  sensitive   = true
}

output "ca_certificate" {
  value       = module.gke.ca_certificate
  description = "CA certificates bundle for the GKE controller"
}
