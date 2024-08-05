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
    google = {
      source  = "hashicorp/google"
      version = ">= 5.38.0"
    }
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = ">= 17.1.0"
    }
  }

  required_version = ">= 1.7.0"
}

provider "gitlab" {
  # remember to export GITLAB_TOKEN environment variable
}
