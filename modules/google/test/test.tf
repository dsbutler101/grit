###############
# TEST MODULE #
###############

module "test-module" {
  source = "./internal"

  manager_service   = var.manager_service
  fleeting_service  = var.fleeting_service
  gitlab_project_id = var.gitlab_project_id
  runner_token      = var.runner_token
  name              = var.name
}

