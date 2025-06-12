terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 6.30.0"
    }
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = ">= 17.0.0"
    }
  }
  backend "http" {}
}

provider "gitlab" {}

provider "google" {
  project = var.google_project
  region  = var.google_region
  zone    = var.google_zone
}

module "blue" {
  source = "../common"

  name              = var.name
  gitlab_project_id = var.gitlab_project_id
  runner_tag_list = [
    var.runner_tag,
    "blue",
  ]
  runner_version_skew = var.runner_version_skew

  google_project = var.google_project
  google_region  = var.google_region
  google_zone    = var.google_zone

  concurrent = 1
}
