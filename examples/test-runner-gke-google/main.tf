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

variable "google_region" {
  type = string
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
  metadata = {
    name        = var.name
    labels      = var.labels
    min_support = "experimental"
  }

  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnetwork_ids[local.metadata.name]
  }
}

provider "google" {
  project = var.google_project
}


module "vpc" {
  source = "../../modules/google/vpc/test/"

  metadata = local.metadata

  google_region = var.google_region

  subnetworks = {
    "${local.metadata.name}" : "10.0.0.0/10"
  }
}

module "cluster" {
  source = "../../modules/google/gke/test/"

  metadata = local.metadata

  # TODO
  google_region = "unused-but-mandatory"

  google_zone = var.google_zone
  nodes_count = 1

  vpc = local.vpc
}

provider "kubectl" {
  host                   = module.cluster.host
  cluster_ca_certificate = module.cluster.ca_certificate
  token                  = module.cluster.access_token
  load_config_file       = false
}

module "operator" {
  source = "../../modules/k8s/operator/test/"

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
  source = "../../modules/k8s/runner/test/"

  metadata  = local.metadata
  namespace = module.operator.namespace
  gitlab    = module.gitlab
}
