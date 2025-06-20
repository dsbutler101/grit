terraform {
  backend "http" {
    address        = "https://gitlab.com/api/v4/projects/48756626/terraform/state/default"
    lock_address   = "https://gitlab.com/api/v4/projects/48756626/terraform/state/default/lock"
    unlock_address = "https://gitlab.com/api/v4/projects/48756626/terraform/state/default/lock"
    lock_method    = "POST"
    unlock_method  = "DELETE"
    retry_wait_min = "5"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.61.0"
    }
    google = {
      source  = "hashicorp/google"
      version = ">= 5.38.0"
    }
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = ">= 17.1.0"
    }
  }

  required_version = ">= 0.14"
}

provider "gitlab" {
  # remember to export GITLAB_TOKEN environment variable
}

provider "aws" {
  region = local.aws_region
}

provider "google" {
  project = local.google_project
  region  = local.google_region
  zone    = local.google_zone
}
