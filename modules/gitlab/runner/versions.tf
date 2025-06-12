terraform {
  required_providers {
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = ">= 16.3.0"
    }
  }
  required_version = ">= 0.14"
}