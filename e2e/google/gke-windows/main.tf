terraform {
  required_providers {
    kubectl = {
      source  = "alekc/kubectl"
      version = "~> 2.0"
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

  backend "http" {}
}

# provider defaults using env vars (GOOGLE_PROJECT etc)
provider "google" {}

data "google_client_config" "current" {}

provider "kubectl" {
  host                   = module.gke_runner.cluster_host
  cluster_ca_certificate = module.gke_runner.cluster_ca_certificate
  token                  = module.gke_runner.cluster_access_token
  load_config_file       = false
}

module "gke_runner" {
  source = "../../../scenarios/google/gke/operator"

  google_region       = data.google_client_config.current.region
  google_zone         = data.google_client_config.current.zone
  subnet_cidr         = "10.0.0.0/16"
  deletion_protection = false
  autoscaling         = { enabled = false, autoscaling_profile = "", auto_provisioning_locations = [], resource_limits = [] }
  labels = {
    "gitlab-project-id" = var.gitlab_project_id
    "e2e"               = "gke-windows-node"
  }

  node_pools = {
    windows = {
      node_count = 1
      node_config = {
        image_type   = "windows_ltsc_containerd"
        machine_type = "n2d-standard-4"
      }
    }
    linux = {
      # 2 x e2-medium required to schedule runner + system pods
      node_count = 2
      node_config = {
        image_type   = "cos_containerd"
        machine_type = "n2d-standard-2"
      }
    }
  }
  name = var.name

  gitlab_project_id  = var.gitlab_project_id
  runner_description = var.name

  runner_tags     = [var.runner_tag]
  config_template = <<EOF
  [[runners]]
    name = ""
    url = "https://gitlab.com/"
    executor = "kubernetes"
    environment = [ "FF_USE_POWERSHELL_PATH_RESOLVER=true" ]
    shell = "powershell"
    [runners.kubernetes]
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
  helper_image = "registry.gitlab.com/gitlab-org/gitlab-runner/gitlab-runner-helper:x86_64-latest-servercore1809"
  // TODO: The default runner image in operator is now registry.gitlab.com/gitlab-org/gitlab-runner:alpine-bleeding which
  // isn't ideal as it doesn't run in the security context of the Operator by default
  runner_image                = "registry.gitlab.com/gitlab-org/ci-cd/gitlab-runner-ubi-images/gitlab-runner-ocp:v${var.runner_version}"
  override_operator_manifests = "file://../../../examples/test-runner-gke-google/operator.k8s.yaml"
}
