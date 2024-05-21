---
stage: Verify
group: Runner
info: To determine the technical writer assigned to the Stage/Group associated with this page, see https://handbook.gitlab.com/handbook/product/ux/technical-writing/#assignments
---

# AWS - Docker autoscaler default scenario

This scenario allows you to spin up GitLab runners with the docker autoscaler executor in AWS.

## How to run this scenario

In order to run this scenario, you need to do 

### Basic usage

In the following example, we will create a runner stack using the default values set in the `variables.tf` file. The only thing you need to edit in set is the `gitlab_project_id`

```tf
terraform {
  required_providers {
    gitlab = {
      source = "gitlabhq/gitlab"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
  ## Uncomment if you want do not want to use GitLab-managed terraform state https://docs.gitlab.com/ee/user/infrastructure/iac/terraform_state.html
  # backend "http" {
  # }
}

module "runner-deployment" {
  source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/aws/linux/docker-autoscaler-default?ref=aws-docker-autoscaler-scenario"

  gitlab_project_id = "YOUR_PROJECT_ID"
}
```

You can either apply this terraform configuration to your AWS account locally using the AWS cli, or use a GitLab CI/CD for it as described in this [Blog Post]()

### Advanced usage

In the advanced example, we are overwriting some of the default values to control things like the autoscaling behaviour of runner, or bring your own VPC. You can check `variables.tf` file, for other variables you can overwrite.

```tf
terraform {
  required_providers {
    gitlab = {
      source = "gitlabhq/gitlab"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
  # Comment if you want do not want to use GitLab-managed terraform state https://docs.gitlab.com/ee/user/infrastructure/iac/terraform_state.html
  backend "http" {
  }
}

module "runner-deployment" {
  source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/aws/linux/docker-autoscaler-default?ref=aws-docker-autoscaler-scenario"

  gitlab_project_id = "YOUR_PROJECT_ID"

  #TODO: Add more configuration
}
```
