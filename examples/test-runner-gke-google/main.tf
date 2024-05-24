terraform {
  # required_version = ">= 0.13"

  required_providers {
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
    google = {
      source  = "hashicorp/google"
      version = ">= 5.30.0"
    }
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = ">= 17.0.0"
    }
  }
}

variable "google_zone" {
  type = string
  # default = "europe-north1-c"
}

variable "google_project" {
  type = string
  # default = "hhoerl-e4e9f672"
}

variable "labels" {
  type    = map(string)
  default = {}
}

variable "name" {
  type = string
}

variable "gitlab_pat" {
  type = string
}

variable "gitlab_project_id" {
  type = string
}

locals {
  # we only have the zone (europe-north1-c), so we split of the last part to get the region (europe-north1)
  region = replace(var.google_zone, "/-[^-]+$/", "")

  # TODO: use the vpc module
  default_vpc_subnetwork = "projects/${var.google_project}/regions/${local.region}/subnetworks/default"
  default_vpc_network    = "projects/${var.google_project}/global/networks/default"

  metadata = {
    name        = var.name
    labels      = var.labels
    min_support = "experimental"
  }
}

provider "google" {
  project = var.google_project
}

module "cluster" {
  source = "../../modules/google/gke/test/"

  metadata = local.metadata

  # TODO
  google_region = "unused-but-mandatory"

  google_zone = var.google_zone
  nodes_count = 1

  vpc = {
    id        = local.default_vpc_network
    subnet_id = local.default_vpc_subnetwork
  }
}

provider "kubectl" {
  host                   = module.cluster.host
  cluster_ca_certificate = module.cluster.ca_certificate
  token                  = module.cluster.access_token
  load_config_file       = false
}

module "operator" {
  source = "../../modules/operator/internal/"

  depends_on = [
    module.cluster
  ]
}

provider "gitlab" {
  token = var.gitlab_pat
}

module "gitlab" {
  source = "../../modules/gitlab/test/"

  metadata           = local.metadata
  project_id         = var.gitlab_project_id
  runner_description = "this is just some runner, brought to you by GRIT"
}

module "runner" {
  source = "../../modules/k8s/runner/internal/"

  name      = "some-runner"
  namespace = module.operator.namespace
  token     = module.gitlab.runner_token
  url       = module.gitlab.url
}
