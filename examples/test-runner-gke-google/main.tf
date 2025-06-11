locals {
  name = "test-gke-google"
}

module "gke_runner" {
  source = "../../scenarios/google/gke/operator/"
  # source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/google/gke/operator"

  google_region = var.google_region
  google_zone   = var.google_zone
  name          = local.name
  node_pools = {
    defaut = {}
  }

  runners = {
    main = {
      runner_token = var.runner_token
    }
  }

  operator = {
    version            = "latest"
    override_manifests = var.override_operator_manifests
  }
}
