# GitLab managed tf state backend

terraform {
  required_providers {
    gitlab = {
      source = "gitlabhq/gitlab"
    }
  }
  backend "http" {
  }
}