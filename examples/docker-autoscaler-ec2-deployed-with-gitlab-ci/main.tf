terraform {
  # Uses the gitlab terraform provider to manage terraform
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
  # Uncomment if you want to use GitLab-managed Terraform state - reccomended if deployed with GitLab CI
  # Documentation: https://docs.gitlab.com/ee/administration/terraform_state.html
  # backend "http" {
  # }
}

locals {
  aws_zone = "us-east-1b"
}

# Valid project id
# How to get the project id? https://docs.gitlab.com/ee/user/project/working_with_projects.html#access-the-project-overview-page-by-using-the-project-id
variable "gitlab_project_id" {
  type      = string
  sensitive = false
}

module "runner-deployment" {
  # Pointing to GRIT's AWS Docker Autoscaler Scenario
  # For more scenarios, see: https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit/-/tree/main/scenarios/
  
  # If you add GRIT to your repository
  source = "../../scenarios/aws/linux/docker-autoscaler-default"

  # Uncomment if you instead want to use remote source
  # source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/aws/linux/docker-autoscaler-default"

  gitlab_project_id = var.gitlab_project_id
}
