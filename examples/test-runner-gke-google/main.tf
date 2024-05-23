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
