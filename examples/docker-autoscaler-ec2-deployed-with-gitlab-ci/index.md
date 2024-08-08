# Docker Autoscaler EC2 Example

Using this example, we use the `aws/linux/docker-autoscaler-default` scenario to deploy an autoscaling runner group to AWS EC2.
It also uses GitLab OpenTofu CI/CD Component to build, test and deploy the runner infrastructure with a GitLab-managed terraform state.

## Getting Started

### Prepare AWS Account

1. Create an AWS account or use existing one
1. Create an IAM user with programmatic access (we provide an example in `example-aws-iam-policy.yml`)
1. Create AWS access key (required to create the runner infrastructure)

**Disclaimer**

This IAM policy is provided as an example and usage comes at your own risk. The author and contributors are not responsible for any security issues arising from the use or misuse of this policy.

### Create GitLab Project

1. Create GitLab Project
1. Create a GitLab [project access token](https://docs.gitlab.com/ee/user/project/settings/project_access_tokens.html#create-a-project-access-token) or [personal access token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html#create-a-personal-access-token) (required to create the runner token)
1. Set CI/CD variables for `GITLAB_TOKEN`, `AWS_DEFAULT_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`

### Add `main.tf` file

The GRIT configuration is used from [`aws/linux/docker-autoscaler-default`](../../scenarios/aws/linux/docker-autoscaler-default/) scenario in the "runner-deployment" module as such `source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/aws/linux/docker-autoscaler-default"`

Steps:

1. Add `main.tf` file
1. Copy content from `./main.tf`
1. Change the `gitlab_project_id` to your project

### Add `gitlab-ci.yml` file

We are now adding a GitLab CI configuration, to build, test, and deploy the runner infrastructure using a GitLab-managed state file.

Steps:

1. Add `.gitlab-ci.yml` file
1. Copy content from `./.gitlab-ci.yml`
1. Make sure there is a Runner available in your project, which can execute the jobs.

### Commit changes & Test

Once you have commited your changes, your CI/CD pipeline should be created based off the `.gitlab-ci.yml` configuration, and allow you to deplot your first `GRIT` runner.
