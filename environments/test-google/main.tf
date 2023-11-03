module "test" {
  source            = "../../modules/google/test"
  manager_provider  = "helm"
  fleeting_service  = "gke"
  gitlab_project_id = var.gitlab_project_id
}