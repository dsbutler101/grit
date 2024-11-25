---
stage: Verify
group: Runner
info: To determine the technical writer assigned to the Stage/Group associated with this page, see https://handbook.gitlab.com/handbook/product/ux/technical-writing/#assignments
---

# Deploy GitLab Runner Manager on AWS

This document describes how to deploy GitLab Runner Manager on AWS by using the Docker Autoscaler executor.

Prerequisites:

- AWS account with administrative access
- AWS CLI installed and configured
- Terraform (version 1.5 or later)
- GitLab account with a project where you have the Maintainer or Owner role
- GitLab personal access token with `api` scope (set as `GITLAB_TOKEN` environment variable)

To deploy GitLab Runner Manager on AWS by using Docker Autoscaler:

1. Configure AWS credentials by using either:

   - AWS CLI:

     ```shell
     aws configure
     ```

   - Environment variables:

     ```shell
     export AWS_ACCESS_KEY_ID="your_access_key"
     export AWS_SECRET_ACCESS_KEY="your_secret_key"
     export GITLAB_TOKEN="your_gitlab_token"  # Required for GitLab provider
     ```

1. Create a new directory for your deployment:

   ```shell
   mkdir my-runner-deployment
   cd my-runner-deployment
   ```

1. Create a `main.tf` file with the following configuration:

   ```hcl
   terraform {
     required_providers {
       gitlab = {
         source  = "gitlabhq/gitlab"
         version = ">=17.0.0"
       }
       aws = {
         source  = "hashicorp/aws"
         version = "~> 5.0"
       }
     }
   }

   provider "aws" {
     region = "us-east-1"  # Change to your preferred region
   }

   provider "gitlab" {
     # Token can be provided via GITLAB_TOKEN environment variable
     # or uncomment and set token here:
     # token = "glpat-xxxxxxxxxxxxxxxxx"
   }

   module "runner-deployment" {
     source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/aws/linux/docker-autoscaler-default"

     # Required: Your GitLab project ID
     gitlab_project_id = "YOUR_PROJECT_ID"

     # Optional configurations with defaults
     aws_region = "us-east-1"
     aws_zone   = "us-east-1a"

     runner_description = "AWS Docker Autoscaler Runner"
     runner_tags        = ["aws", "docker", "autoscaler"]

     max_instances = 10
     ephemeral_runner = {
       machine_type = "t3.micro"
       source_image = ""  # Leave empty for default AMI
     }

     autoscaling_policy = {
       scale_min    = 0
       scale_factor = 0.7
     }

     capacity_per_instance = 10
     gitlab_url            = "https://gitlab.com"  # Change for self-hosted GitLab
   }
   ```

1. Deploy the runner:

   ```shell
   # Initialize Terraform
   terraform init

   # Preview changes
   terraform plan

   # Apply configuration
   terraform apply
   ```

## Configuration Options

The following are the key configuration parameters:

- AWS Settings:
  - `aws_region`: Your AWS region (for example, `us-east-1`)
  - `aws_zone`: Availability zone (for example, `us-east-1a`)
  - `ephemeral_runner.machine_type`: EC2 instance type

- Runner Settings:
  - `gitlab_project_id`: Your GitLab project ID
  - `runner_tags`: Tags for job matching
  - `max_instances`: Maximum number of EC2 instances
  - `capacity_per_instance`: Concurrent jobs per instance

- Autoscaling:
  - `autoscaling_policy.scale_min`: Minimum instances
  - `autoscaling_policy.scale_factor`: Idle threshold (between `0.0` to `1.0`)
