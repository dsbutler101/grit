###############
# TEST MODULE #
###############

module "test-module" {
  source = "./internal"

  fleeting_service   = var.fleeting_service
  gitlab_project_id  = var.gitlab_project_id
  manager_provider   = var.manager_provider
  gitlab_runner_tags = var.gitlab_runner_tags
}