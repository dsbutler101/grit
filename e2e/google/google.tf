variable "runner_token" {}
variable "name" {}
variable "job_id" {}
variable "google_region" {}
variable "google_zone" {}
variable "google_project" {}

output "instance_group_name" {
  value = module.fleeting.instance_group_name
}

locals {
  metadata = {
    name = var.name
    labels = tomap({
      job_id = var.job_id
      env    = "grit-e2e"
    })
    min_support = "experimental"
  }
}

module "iam" {
  source   = "../../modules/google/iam"
  metadata = local.metadata
}

module "vpc" {
  source   = "../../modules/google/vpc/prod"
  metadata = local.metadata

  google_region = var.google_region

  subnetworks = {
    "${var.name}-runner-manager"    = "10.0.0.0/29"
    "${var.name}-ephemeral-runners" = "10.1.0.0/21"
  }
}

module "fleeting" {
  source   = "../../modules/google/fleeting/prod"
  metadata = local.metadata
  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnetwork_ids["${var.name}-ephemeral-runners"]
  }

  fleeting_service      = "gce"
  google_project        = var.google_project
  google_zone           = var.google_zone
  service_account_email = module.iam.service_account_email
  machine_type          = "n2d-standard-2"
  manager_subnet_cidr   = module.vpc.subnetwork_cidrs["${var.name}-runner-manager"]

}

module "runner" {
  source   = "../../modules/google/runner/prod"
  metadata = local.metadata

  google_project = var.google_project
  google_zone    = var.google_zone

  service_account_email = module.iam.service_account_email

  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnetwork_ids["${var.name}-runner-manager"]
  }

  gitlab_url   = local.gitlab.url
  runner_token = local.gitlab.runner_token

  executor = "docker-autoscaler"

  fleeting_instance_group_name = module.fleeting.instance_group_name

  machine_type = "n2d-standard-2"
}

locals {
  gitlab = {
    runner_token = var.runner_token
    url          = "https://gitlab.com"
  }
}
