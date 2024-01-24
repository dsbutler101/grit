#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

######################
# GITLAB TEST MODULE #
######################

module "gitlab" {
  source             = "../internal"
  name               = var.metadata.name
  project_id         = var.project_id
  url                = var.url
  runner_description = var.runner_description
  runner_tags        = var.runner_tags
}
