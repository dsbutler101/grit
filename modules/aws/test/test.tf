###############
# TEST MODULE #
###############

module "test-module" {
  source = "./internal"

  fleeting_service = var.fleeting_service
  manager_service  = var.manager_service

  runner_token       = var.runner_token
  gitlab_project_id  = var.gitlab_project_id
  gitlab_runner_tags = var.gitlab_runner_tags
  name               = var.name
}

