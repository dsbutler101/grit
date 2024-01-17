output "gke_cluster_name" {
  value = try(module.gke-cluster[0].name, "")
}

output "gke_cluster_host" {
  value = try(module.gke-cluster[0].host, "")
}

output "gke_cluster_access_token" {
  value     = try(module.gke-cluster[0].access_token, "")
  sensitive = true
}

output "gke_cluster_ca_certificate" {
  value = try(module.gke-cluster[0].ca_certificate, "")
}

