output "cluster_host" {
  description = "The GKE cluster's control plane URL"
  value       = module.cluster.host
}

output "cluster_ca_certificate" {
  description = "The GKE cluster's CA certificate"
  value       = module.cluster.ca_certificate
}

output "cluster_access_token" {
  description = "The GKE cluster's admin token"
  value       = module.cluster.access_token
  sensitive   = true
}

output "supported_operator_versions" {
  value = module.operator.supported_operator_versions
}
