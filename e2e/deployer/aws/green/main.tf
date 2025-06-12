terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.61.0"
    }
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = ">= 17.0.0"
    }
  }
  backend "http" {}
}

provider "gitlab" {}

provider "aws" {}

module "green" {
  source = "../common"

  name              = var.name
  gitlab_project_id = var.gitlab_project_id
  runner_tag_list = [
    var.runner_tag,
    "green",
  ]
  runner_version_skew = var.runner_version_skew
}
