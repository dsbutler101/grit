module "test" {
  source            = "../../modules/test"
  manager_provider  = "helm"
  capacity_provider = "gke"
  gitlab_project_id = var.gitlab_project_id
}