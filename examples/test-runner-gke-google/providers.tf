provider "google" {
  project = var.google_project
}

provider "gitlab" {
  token = var.gitlab_pat
}

provider "kubectl" {
  host                   = module.gke_runner.cluster_host
  cluster_ca_certificate = module.gke_runner.cluster_ca_certificate
  token                  = module.gke_runner.cluster_access_token
  load_config_file       = false
}
