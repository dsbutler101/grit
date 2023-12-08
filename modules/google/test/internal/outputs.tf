output "gke_cluster_name" {
  value = module.gke-cluster[0].name
}

output "gke_cluster_host" {
  value = module.gke-cluster[0].host
}

output "gke_cluster_access_token" {
  value = module.gke-cluster[0].access_token
}

output "gke_cluster_ca_certificate" {
  value = module.gke-cluster[0].ca_certificate
}

