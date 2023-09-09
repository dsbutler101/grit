provider "google" {
  project = var.project_id
  zone    = var.zone
}

provider "helm" {
  kubernetes {
    host                   = google_container_cluster.primary.endpoint
    token                  = data.google_client_config.provider.access_token
    cluster_ca_certificate = base64decode(google_container_cluster.primary.master_auth[0].cluster_ca_certificate)
  }
}

provider "gitlab" {
  token = var.gitlab_token
}
