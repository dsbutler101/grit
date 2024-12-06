terraform {
  required_providers {
    kubectl = {
      source = "alekc/kubectl"
    }
  }
}

variable "gitlab_runner_token" {}
variable "google_project" {}
variable "google_region" {}
variable "google_zone" {}
variable "name" {}

provider "google" {
  project = var.google_project
  region  = var.google_region
}

provider "kubectl" {
  host                   = module.cluster.host
  token                  = module.cluster.access_token
  cluster_ca_certificate = module.cluster.ca_certificate
  load_config_file       = false
}

locals {
  metadata = {
    name        = var.name
    labels      = {}
    min_support = "experimental"
  }
  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnetwork_ids[local.metadata.name]
  }
  gitlab_credentials = {
    url          = "https://gitlab.com"
    runner_token = var.gitlab_runner_token
  }
}

module "vpc" {
  source        = "../../../modules/google/vpc/prod/"
  metadata      = local.metadata
  google_region = var.google_region
  subnetworks = zipmap(
    [local.metadata.name],
    ["10.0.0.0/24"],
  )
}

module "cluster" {
  source        = "../../../modules/google/gke/prod/"
  metadata      = local.metadata
  google_region = var.google_region
  google_zone   = var.google_zone

  deletion_protection = false
  node_pools = {
    windows = {
      node_count = 1
      node_config = {
        image_type = "windows_ltsc_containerd"
      }
    }
    linux = {
      node_count = 1
      node_config = {
        image_type   = "cos_containerd"
        machine_type = "e2-medium"
      }
    }
  }
  vpc = local.vpc
}

module "operator" {
  source     = "../../../modules/k8s/operator/prod/"
  metadata   = local.metadata
  depends_on = [module.cluster]
}

module "runner" {
  source          = "../../../modules/k8s/runner/prod/"
  metadata        = local.metadata
  namespace       = module.operator.namespace
  helper_image    = "gitlab/gitlab-runner-helper:x86_64-latest-servercore1809"
  gitlab          = local.gitlab_credentials
  config_template = <<EOF
[[runners]]
  name = ""
  url = "https://gitlab.com/"
  executor = "kubernetes"
  environment = [ "FF_USE_POWERSHELL_PATH_RESOLVER=true" ]
  shell = "powershell"
  [runners.kubernetes]
    image = "mcr.microsoft.com/powershell:lts-windowsservercore-1809"
    pod_labels_overwrite_allowed = ""
    service_account_overwrite_allowed = ""
    pod_annotations_overwrite_allowed = ""
    [runners.kubernetes.node_selector]
      "kubernetes.io/arch" = "amd64"
      "kubernetes.io/os" = "windows"
      "node.kubernetes.io/windows-build" = "10.0.17763"
    [runners.kubernetes.pod_security_context]
    [runners.kubernetes.volumes]
    [runners.kubernetes.dns_config]
EOF
  envvars = {
    "KUBERNETES_POLL_TIMEOUT" = "3600"
    "FF_TIMESTAMPS"           = "true"
  }
}

output "host" {
  value       = module.cluster.host
  description = "Host of the GKE controller"
}

output "access_token" {
  value       = module.cluster.access_token
  description = "Access token for the GKE controller"
  sensitive   = true
}

output "ca_certificate" {
  value       = module.cluster.ca_certificate
  description = "CA certificates bundle for the GKE controller"
  sensitive   = true
}
