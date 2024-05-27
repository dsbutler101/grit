provider "google" {
  project = var.google_project
}

provider "gitlab" {
  token = var.gitlab_pat
}

provider "kubectl" {
  host                   = module.cluster.host
  cluster_ca_certificate = module.cluster.ca_certificate
  token                  = module.cluster.access_token
  load_config_file       = false
}
