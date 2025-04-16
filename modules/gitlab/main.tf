#######################
# METADATA VALIDATION #
#######################

module "validate_name" {
  source = "../internal/validation/name"
  name   = var.metadata.name
}

module "validate_support" {
  source   = "../internal/validation/support"
  use_case = "any"
  use_case_support = tomap({
    "any" = "experimental"
  })
  min_support = var.metadata.min_support
}

resource "gitlab_user_runner" "primary" {
  description = "${var.runner_description} ${var.metadata.name}_GRIT"
  runner_type = var.runner_type
  tag_list    = var.runner_tags
  untagged    = length(var.runner_tags) == 0 ? true : false
  group_id    = var.runner_type == "group_type" ? var.group_id : null
  project_id  = var.runner_type == "project_type" ? var.project_id : null

  // For non-project runners value doesn't matter
  locked = var.runner_type == "project_type" ? var.lock_project_runner : true

  lifecycle {
    precondition {
      condition     = !(var.runner_type == "group_type" && var.group_id == "")
      error_message = "The group_id must be set when runner_type is 'group_type'."
    }

    precondition {
      condition     = !(var.runner_type == "project_type" && var.project_id == "")
      error_message = "The project_id must be set when runner_type is 'project_type'."
    }
  }
}
