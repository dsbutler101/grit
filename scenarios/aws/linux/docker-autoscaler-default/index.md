# Scenario: AWS - Linux - Docker Autoscaler default

This scenario template deploys GitLab Runner to AWS, with configuration
that supports autoscaling for Linux through the `docker-autoscaler` executor.

## Prerequisites

To use this scenario, you must meet the following prerequisites:

### Terraform and AWS setup

To use this scenario, you must have:

- [Terraform prerequisites](../../../index.md#prerequisites)
- [AWS prerequisites](../../../index.md)
- Terraform 1.5 or later

### IAM Role

In order to control AWS resources, an AWS key must be provided by setting the following environment variables:

```markdown
export AWS_SECRET_ACCESS_KEY=<value>
export AWS_ACCESS_KEY_ID=<value>
```

Below, the list of AWS services used and the recommended IAM policy.

#### AWS Services

- Amazon CloudWatch
- Amazon EC2
- Amazon EC2 Auto Scaling
- AWS Identity and Access Management
- AWS Security Token Service

#### IAM Policy

User's responsibility: 

- Review, test, and modify this policy to align with your company's security policies and requirements
- Minimum permissions: Tailor the policy to grant only necessary permissions for your specific use case
- Regular maintenance: Conduct periodic audits and updates to maintain security as your needs evolve
- Best practices: Consult your security team and AWS documentation for current IAM policy best practices

Always apply the principle of least privilege when working with IAM policies. We provide an example in [IAM Policy](../../../../examples/docker-autoscaler-ec2-deployed-with-gitlab-ci/index.md#prepare-aws-account)] which can be used as a reference.

## Variables

You can use variables to control the behavior of the scenario.

Variables can be:

- **Required**: Variables must be provided when you define the module and do not have
  a default value.

- **Not required with a default value**: Variables are required for the scenario to work properly, but you
  can use the provided default values to experiment with the scenario.

- **Not required with no default value**: Variables are optional and don't need to be provided
  unless a specific use case requires changes in the default configuration.

- **Simple**: Variables use simple types as `string`, `number` or `boolean`.

- **Complex**: Variables are either lists, maps, or objects, or combination of these types.

| Name                     | Type                                                     | Required? | Default value | Description                                                                                                                                                                                         |
| ------------------------ | -------------------------------------------------------- | --------- | ------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `name`                   | `string`                                                 | yes       |               | Name of the deployment. Must be unique in scope of an AWS Account (20 chars max) project.                                                                                                           |
| `labels`                 | `map(string)`                                            | no        |               | Arbitrary list of `key=value` pairs that are added as labels to resources created by GRIT.                                                                                                          |
| `aws_region`             | `string`                                                 | yes       |               | AWS region chosen for the deployment.                                                                                                                                                               |
| `aws_zone`               | `string`                                                 | yes       |               | AWS availability zone chosen for the deployment.                                                                                                                                                    |
| `gitlab_url`             | `string`                                                 | yes       |               | URL of GitLab instance.                                                                                                                                                                             |
| `runner_token`           | `string`                                                 | yes       |               | Authentication token of the runner to deploy. See [how to obtain the token](https://docs.gitlab.com/ee/ci/runners/runners_scope.html#create-an-instance-runner-with-a-runner-authentication-token). |
| `runner_machine_type`    | `string`                                                 | no        |               | Machine type for the runner manager instance. If not provided, GRIT uses one of predefined choices based on the value of defined concurrency.                                                       |
| `concurrent`             | `number`                                                 | no        | 1             | Value for `config.toml`'s `concurrent` setting. Defines maximum total number of jobs executed concurrently by the runner.                                                                           |
| `runners_global_section` | `string`                                                 | no        |               | [Allows to customize](#runners_global_section-customization) the global part of `[[runners]]` section in generated `config.toml`.                                                                   |
| `runners_docker_section` | `string`                                                 | no        |               | [Allows to customize](#runners_docker_section-customization) the global part of `[runners.docker]` section in generated `config.toml`.                                                              |
| `capacity_per_instance`  | `number`                                                 | no        | 1             | Maximum number of jobs to be executed concurrently on one autoscaled ephemeral VM.                                                                                                                  |
| `max_instances`          | `number`                                                 | no        | 200           | Maximum number of ephemeral instances (in all possible states) that runner maintains.                                                                                                               |
| `max_use_count`          | `number`                                                 | no        | 1             | Maximum number of jobs executed on a single autoscaled ephemeral VM, after which the VM is marked for deletion.                                                                                     |
| `autoscaling_policies`   | [`list(object)`](#autoscaling_policies-object-structure) | no        |               | List of objects defining autoscaling policies.                                                                                                                                                      |
| `ephemeral_runner`       | [`object`](#ephemeral_runner-object-structure)           | no        |               | Configuration of autoscaled ephemeral VM.                                                                                                                                                           |

### `autoscaling_policy` object structure

| Key                  | Type     | Description                                                                                                                                                                                                                                 |
| -------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `scale_min`          | `number` | The minimum size of the autoscaling instances fleet. This defines the number of `idle` Taskscaler slots and therefore Fleeting instances that should be sustained for all time to have space for jobs execution. Integer value is expected. |
| `idle_time`          | `string` | Minimal duration for which the `idle` instances should running, even if not used. Uses Go's time format, for example `1h20m30s`.                                                                                                            |
| `scale_factor`       | `number` | If used, the number of idle slots are calculated as `scale_factor * in_use_slots`, but not less than defined with `scale_min`. A `float64` value higher than 0 is expected.                                                                 |
| `scale_factor_limit` | `number` | Usable only when `scale_factor` is in use. If defined, the maximum value of `idle` is calculated with the equation described for `scale_factor`. An integer value is expected.                                                              |

If not defined by the user, the following default is applied:

```terraform
object {
  scale_min          = 1
  scale_factor       = 0
}
```

### `ephemeral_runner` object structure

| Key            | Type     | Description                                                                                         |
| -------------- | -------- | --------------------------------------------------------------------------------------------------- |
| `disk_type`    | `string` | [Volume type](https://aws.amazon.com/ebs/volume-types/) to be used for autoscaled ephemeral VMs.    |
| `disk_size`    | `number` | Disk size (in GiB) to be used for autoscaled ephemeral VMs. Integer value is expected.              |
| `machine_type` |          | [Machine type](https://aws.amazon.com/ec2/instance-types/) to be used for autoscaled ephemeral VMs. |
| `source_image` |          | Source AMI from which autoscaled ephemeral VMs are started.                                         |

If not defined by the user, the following default is applied:

```terraform
object {
  disk_type    = "gp3"
  disk_size    = 25
  machine_type = "t3.medium"
  source_image = module.ami_lookup.ami_id
}
```

## Usage

### Terraform code

Consider the following examples of the Terraform code that uses this scenario:

- Simple: defines only the required variables of the module.
- Advanced: uses all available variables.

Both examples use the same common part that contains:

- Definition of version requirements (for Terraform itself and used providers).
- Initialization of local and external variables.
- Initialization of the `aws` provider.

In both examples, the `runner_token` variable is an external variable marked
as `sensitive`.

The `runner_token` defines the identity of the runner. It provides access to all projects
and jobs for the runner, and must be treated as a secret value.

You should not put the `runner_token` variable into the `*.tf` file
or commit it to any version control system. Instead, it should be passed
to Terraform as a variable, as an environment variable, with the CLI flag or
a variables file. Read more about [how to use input variables with Terraform](https://developer.hashicorp.com/terraform/language/values/variables#assigning-values-to-root-module-variables).

#### Simple example, deployed using GitLab CI

To quickest way to get started is to follow the [Docker Autoscaler EC2 - deployed with GitLab CI](../../../../examples/docker-autoscaler-ec2-deployed-with-gitlab-ci) example. This will get you a semi-production ready autoscaling GitLab Runner, deployed using GitLab CI with a remote Terraform state file, in under 10 mins.

#### Simple example with local deployment

To see how to apply this example locally with Terraform CLI, see [deploy locally with terraform CLI](../../../index.md#deploy-locally-with-terraform-cli).

```terraform
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
    # Uncomment this if you want GitLab to manage your Terraform state
    # Documentation: https://docs.gitlab.com/ee/administration/terraform_state.html
    #backend "http" {
    #}
  }
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
  source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/aws/linux/docker-autoscaler-default"

  gitlab_project_id = var.gitlab_project_id
}
```

#### Advanced example

```terraform
terraform {
  required_providers {
    gitlab = {
      source = "gitlabhq/gitlab"
      version = ">=17.0.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
}

provider "gitlab" {
  base_url = var.gitlab_url
}

locals {
  name       = "gritexample"
  aws_zone   = "us-east-2a"
  aws_region = "us-east-2"
}

variable "runner_token" {
  type      = string
  sensitive = true
}

variable "gitlab_url" {
  type      = string
  sensitive = false
  default   = "https://gitlab.com"
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
  source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/aws/linux/docker-autoscaler-default?ref=aws-docker-autoscaler-scenario"

  name                  = local.name
  aws_zone              = local.aws_zone
  aws_region            = local.aws_region
  gitlab_url            = var.gitlab_url
  gitlab_project_id     = var.project_id
  runner_token          = var.runner_token
  capacity_per_instance = 2
  max_use_count         = 10
  concurrent            = 1000

  ephemeral_runner = {
    disk_size    = 50
  }

  runners_global_section = <<EOS
  environment = [
    "DOCKER_TLS_CERTDIR=",
  ]
EOS

  runners_docker_section = <<EOS
    volumes = [
      "/certs/client"
    ]

    privileged = true
EOS
}
```
