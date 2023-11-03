terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.74.0"
    }
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = "16.3.0"
    }
  }
  required_version = ">= 0.14"
}
