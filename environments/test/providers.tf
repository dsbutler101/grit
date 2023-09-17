provider "gitlab" {
  base_url = var.gitlab_url
  token    = var.gitlab_token
}

provider "google" {
  project = var.gcp_project_id
  zone    = var.zone
}

provider "helm" {
  kubernetes {
    host                   = module.test.gke-cluster.host
    token                  = module.test.gke-cluster.access_token
    cluster_ca_certificate = base64decode(module.test.gke-cluster.ca_certificate)
  }
}
