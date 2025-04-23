terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.30.0"
    }
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = ">= 17.0.0"
    }
  }

  backend "http" {}
}

locals {
  metadata = {
    name = var.name
    labels = tomap({
      gitlab_project_id = var.gitlab_project_id
      env               = "grit-e2e"
    })
    min_support = "experimental"
  }
}

provider "gitlab" {}

module "gitlab" {
  source             = "../../modules/gitlab"
  metadata           = local.metadata
  url                = "https://gitlab.com"
  project_id         = var.gitlab_project_id
  runner_description = var.name
  runner_tags        = [var.runner_tag]
}

# provider defaults using env vars (GOOGLE_PROJECT etc)
provider "google" {}

data "google_client_config" "current" {}

module "iam" {
  source   = "../../modules/google/iam"
  metadata = local.metadata
}

module "vpc" {
  source   = "../../modules/google/vpc"
  metadata = local.metadata

  google_region = data.google_client_config.current.region

  subnetworks = {
    "${var.name}-runner-manager"    = "10.0.0.0/29"
    "${var.name}-ephemeral-runners" = "10.1.0.0/21"
  }
}

module "fleeting" {
  source   = "../../modules/google/fleeting"
  metadata = local.metadata
  vpc = {
    enabled          = module.vpc.enabled
    id               = module.vpc.id
    subnetwork_ids   = module.vpc.subnetwork_ids
    subnetwork_cidrs = module.vpc.subnetwork_cidrs
  }
  manager_subnet_name = "${var.name}-runner-manager"
  runners_subnet_name = "${var.name}-ephemeral-runners"

  fleeting_service      = "gce"
  google_project        = data.google_client_config.current.project
  subnetwork_project    = data.google_client_config.current.project
  google_zone           = data.google_client_config.current.zone
  service_account_email = module.iam.service_account_email
  machine_type          = "n2d-standard-2"
}

module "runner" {
  source   = "../../modules/google/runner"
  metadata = local.metadata

  google_project     = data.google_client_config.current.project
  subnetwork_project = data.google_client_config.current.project
  google_zone        = data.google_client_config.current.zone

  service_account_email = module.iam.service_account_email

  vpc = {
    enabled          = module.vpc.enabled
    id               = module.vpc.id
    subnetwork_ids   = module.vpc.subnetwork_ids
    subnetwork_cidrs = module.vpc.subnetwork_cidrs
  }
  manager_subnet_name = "${var.name}-runner-manager"

  gitlab_url   = module.gitlab.url
  runner_token = module.gitlab.runner_token

  runner_version = "v${var.runner_version}"

  executor = "docker-autoscaler"

  fleeting_instance_group_name = module.fleeting.instance_group_name

  machine_type = "n2d-standard-2"
}
