module "gke_runner" {
  source = "../../scenarios/google/gke/operator/"
  # source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/google/gke/operator"

  name   = var.name
  labels = var.labels

  google_region      = var.google_region
  google_zone        = var.google_zone
  subnet_cidr        = var.subnet_cidr
  node_count         = var.node_count
  gitlab_pat         = var.gitlab_pat
  gitlab_project_id  = var.gitlab_project_id
  runner_description = var.runner_description
}
