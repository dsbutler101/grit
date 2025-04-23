locals {
  metadata = {
    name        = var.name
    min_support = "none"
    labels = {
      "env" = "grit-e2e"
    }
  }

  subnetworks = {
    "${var.name}-runner-managers"   = "10.0.0.0/29"
    "${var.name}-ephemeral-runners" = "10.1.0.0/21"
  }
}

module "gitlab" {
  source             = "../../../../modules/gitlab"
  metadata           = local.metadata
  url                = "https://gitlab.com"
  project_id         = var.gitlab_project_id
  runner_description = var.name
  runner_tags        = var.runner_tag_list
}

module "vpc" {
  source   = "../../../../modules/google/vpc"
  metadata = local.metadata

  google_region = var.google_region
  subnetworks   = local.subnetworks
}

module "iam" {
  source = "../../../../modules/google/iam"

  metadata = local.metadata
}

module "runner" {
  source   = "../../../../modules/google/runner"
  metadata = local.metadata

  google_project     = var.google_project
  google_zone        = var.google_zone
  subnetwork_project = var.google_project

  service_account_email = module.iam.service_account_email

  vpc = {
    enabled          = module.vpc.enabled
    id               = module.vpc.id
    subnetwork_ids   = module.vpc.subnetwork_ids
    subnetwork_cidrs = module.vpc.subnetwork_cidrs
  }
  manager_subnet_name = "${var.name}-runner-managers"

  gitlab_url     = module.gitlab.url
  runner_token   = module.gitlab.runner_token
  runner_version = var.runner_version

  executor = "shell"

  concurrent     = var.concurrent
  check_interval = 3

  request_concurrency = var.concurrent > 10 ? 10 : var.concurrent

  fleeting_instance_group_name = ""

  runner_wrapper = {
    enabled = true
  }
}
