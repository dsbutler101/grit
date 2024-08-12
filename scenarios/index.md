---
stage: Verify
group: Runner
info: To determine the technical writer assigned to the Stage/Group associated with this page, see https://handbook.gitlab.com/handbook/product/ux/technical-writing/#assignments
---

# Predefined scenario templates

The GitLab Runner Infrastructure Toolkit (GRIT) provides a set of Terraform modules
that can be composed together to create a GitLab Runner deployment.

In addition to components, GRIT provides predefined
scenario templates that compose the lower-level building blocks into
tested, working setups.

Scenario templates provide a higher-level Terraform module
with a limited number of control variables.

## Available scenarios

### Google cloud

1. [Linux - Docker Autoscaler default](google/linux/docker-autoscaler-default)

### AWS

1. [Linux - Docker Autoscaler default](aws/linux/docker-autoscaler-default)

## Prerequisites

### Terraform

GRIT is a library of Terraform modules. To use it, you must have a working Terraform setup.

Depending on your usage, execution of Terraform might typically be done with automation.
For basic experimentation on local machines, an installation of Terraform CLI is
required.

## Using Scenarios

## Deployment

### Deploy using GitLab CI

You can use GitLab CI to deploy a runner. The fastest way to get you started is to use the [OpenTofu GitLab CI/CD component](https://gitlab.com/components/opentofu). See [Docker Autoscaler EC2 - Deployed with GitLab CI](../examples/docker-autoscaler-ec2-deployed-with-gitlab-ci/index.md#add-gitlab-ciyml-file) for an example implementation.

### Deploy locally with Terraform CLI

Read more about [how to install Terraform CLI](https://developer.hashicorp.com/terraform/install).

When the Terraform code is ready, in the directory where the `*.tf` files are created, you can either deploy the infrastructure locally using `Terraform CLI` or 
complete the following steps to execute the Terraform code:

#### Initialize

Use the `init` call to initialize the directory that contains the Terraform configuration files.
The `init` call downloads all providers and external modules referenced by the code.:

```shell
terraform init
```

You must run this command when:

- You first use Terraform in this directory.
- Version definitions in `terraform` block are changed.
- The GRIT code has been updated.

#### Plan

Use the `plan` call to compare the local Terraform state stored in the file with the
resources in AWS.

```shell
terraform plan -var runner_token="your-glrt-runner-token" -out plan.out
```

The `plan` call prints all changes that will be provisioned and stores them
in the `plan.out` file that you use in the last step to execute Terraform.

The value for `runner_token` variable is passed with the flag, `-var runner_token="glrt-runner-token-here"`.
Alternatively, you can export a `TF_VAR_runner_token="glrt-runner-token-here"` before
you run the `plan` call, for example:

```shell
export TF_VAR_runner_token="glrt-runner-token-here"
terraform plan -out plan.out
```

#### Apply the plan

Use the `apply` call to execute runner provisioning steps in the `plan.out` file. Runner
provisioning is executed with calls to the AWS API.

```shell
terraform apply plan.out
```

### Access to the runner manager instance

After the `terraform apply` step is finished, Terraform prints defined outputs.
In the advanced and simple examples, Terraform prints the external IP of the
instance where runner manager was deployed.

This instance can be now accessed with SSH.
