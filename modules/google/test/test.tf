###############
# TEST MODULE #
###############

module "test-module" {
  source = "./internal"

  fleeting_service  = var.fleeting_service
  gitlab_project_id = var.gitlab_project_id
  manager_service   = var.manager_service
  runner_token      = var.runner_token
  name              = ""
}

