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

module "green" {
  source = "../common"

  name              = var.name
  gitlab_project_id = var.gitlab_project_id
  runner_tag_list = [
    var.runner_tag,
    "green",
  ]
  runner_version = var.runner_version

  google_project = var.google_project
  google_region  = var.google_region
  google_zone    = var.google_zone

  concurrent = 2
}
