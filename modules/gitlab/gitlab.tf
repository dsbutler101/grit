#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../internal/validation/support"
  use_case = "any"
  use_case_support = tomap({
    "any" = "experimental"
  })
  min_support = var.metadata.min_support
}

######################
# GITLAB PROD MODULE #
######################

resource "gitlab_user_runner" "primary" {
  description = "${var.runner_description} ${var.metadata.name}_GRIT"
  runner_type = "project_type"
  project_id  = var.project_id
  tag_list    = var.runner_tags
  untagged    = length(var.runner_tags) == 0 ? true : false
}

