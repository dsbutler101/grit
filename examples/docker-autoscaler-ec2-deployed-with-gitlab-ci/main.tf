terraform {
  required_providers {
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = ">=17.0.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
  # Uncomment to use GitLab-managed Terraform state - recommended if deployed with GitLab CI
  # Documentation: https://docs.gitlab.com/ee/administration/terraform_state.html
  # backend "http" {
  # }
}

locals {
  aws_zone = "us-east-1b"
}

variable "gitlab_project_id" {
  type        = string
  description = "The GitLab project ID where the Terraform code is being executed."
}

module "runner-deployment" {
  source = "../../scenarios/aws/linux/docker-autoscaler-default"

  # Uncomment if you instead want to use remote source
  # source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/aws/linux/docker-autoscaler-default"

  gitlab_project_id = var.gitlab_project_id
}
