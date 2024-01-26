#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../../internal/validation/support"
  use_case = "any"
  use_case_support = tomap({
    "any" = "experimental"
  })
  min_support = var.metadata.min_support
}

######################
# GITLAB PROD MODULE #
######################

module "gitlab" {
  source             = "../internal"
  name               = var.metadata.name
  project_id         = var.project_id
  url                = var.url
  runner_description = var.runner_description
  runner_tags        = var.runner_tags
}
