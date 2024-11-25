---
stage: Verify
group: Runner
info: To determine the technical writer assigned to the Stage/Group associated with this page, see https://handbook.gitlab.com/handbook/product/ux/technical-writing/#assignments
---

# Scenario: Google - Linux - Docker Autoscaler default

This scenario template deploys GitLab Runner to Google Cloud, with configuration
that supports autoscaling for Linux through the `docker-autoscaler` executor.

## Prerequisites

To use this scenario, you must meet the following prerequisites:

### Terraform and Google Cloud setup

To use this scenario, you must have:

- [Terraform prerequisites](../../../index.md#prerequisites)
- [Google Cloud prerequisites](../../../index.md#google-cloud-integration-prerequisites)
- Terraform 1.5 or later to use Terraform features and syntax specific to this scenario

### API services

The following API services must be enabled in the Google Cloud project:

- `cloudkms.googleapis.com`
- `compute.googleapis.com`
- `iam.googleapis.com`
- `cloudresourcemanager.googleapis.com`
- `iamcredentials.googleapis.com`
- `oslogin.googleapis.com`

Ask someone with `owner` access to your Google Cloud project to run the following
command:

```shell
gcloud services enable \
    cloudkms.googleapis.com \
    compute.googleapis.com \
    iam.googleapis.com \
    cloudresourcemanager.googleapis.com \
    iamcredentials.googleapis.com \
    oslogin.googleapis.com
```

### Google Cloud permissions for Terraform execution

The actors that execute the Terraform code must have the
following [permissions](https://cloud.google.com/kms/docs/reference/permissions-and-roles) in Google Cloud:

<details>
<summary>Required permissions</summary>

- `cloudkms.cryptoKeyVersions.destroy`
- `cloudkms.cryptoKeyVersions.list`
- `cloudkms.cryptoKeyVersions.useToEncrypt`
- `cloudkms.cryptoKeys.create`
- `cloudkms.cryptoKeys.get`
- `cloudkms.cryptoKeys.update`
- `cloudkms.keyRings.create`
- `cloudkms.keyRings.get`
- `compute.disks.create`
- `compute.firewalls.create`
- `compute.firewalls.delete`
- `compute.firewalls.get`
- `compute.instanceGroupManagers.create`
- `compute.instanceGroupManagers.delete`
- `compute.instanceGroupManagers.get`
- `compute.instanceGroups.create`
- `compute.instanceGroups.delete`
- `compute.instanceTemplates.create`
- `compute.instanceTemplates.delete`
- `compute.instanceTemplates.get`
- `compute.instanceTemplates.useReadOnly`
- `compute.instances.create`
- `compute.instances.delete`
- `compute.instances.get`
- `compute.instances.setLabels`
- `compute.instances.setMetadata`
- `compute.instances.setServiceAccount`
- `compute.instances.setTags`
- `compute.networks.create`
- `compute.networks.delete`
- `compute.networks.get`
- `compute.networks.updatePolicy`
- `compute.regionOperations.get`
- `compute.subnetworks.create`
- `compute.subnetworks.delete`
- `compute.subnetworks.get`
- `compute.subnetworks.use`
- `compute.subnetworks.useExternalIp`
- `compute.zones.get`
- `iam.roles.create`
- `iam.roles.delete`
- `iam.roles.get`
- `iam.roles.list`
- `iam.roles.update`
- `iam.serviceAccounts.actAs`
- `iam.serviceAccounts.create`
- `iam.serviceAccounts.delete`
- `iam.serviceAccounts.get`
- `iam.serviceAccounts.list`
- `resourcemanager.projects.get`
- `resourcemanager.projects.getIamPolicy`
- `resourcemanager.projects.setIamPolicy`
- `storage.buckets.create`
- `storage.buckets.delete`
- `storage.buckets.get`
- `storage.buckets.getIamPolicy`
- `storage.buckets.setIamPolicy`

</details>

You can also create a [custom role](https://cloud.google.com/iam/docs/creating-custom-roles)
with these permissions. You can then assign this role to the user or service account
responsible for provisioning the GRIT Terraform configuration.

Ask someone with `owner` access to your Google Cloud project to run the following
commands:

<details>

```shell
cat > grit-provisioner-role.json <<EOF
{
  "title": "GRITProvisioner",
  "description": "A role with minimum list of permissions required for GRIT provisioning",
  "includedPermissions": [
    "cloudkms.cryptoKeyVersions.destroy",
    "cloudkms.cryptoKeyVersions.list",
    "cloudkms.cryptoKeyVersions.useToEncrypt",
    "cloudkms.cryptoKeys.create",
    "cloudkms.cryptoKeys.get",
    "cloudkms.cryptoKeys.update",
    "cloudkms.keyRings.create",
    "cloudkms.keyRings.get",
    "compute.disks.create",
    "compute.firewalls.create",
    "compute.firewalls.delete",
    "compute.firewalls.get",
    "compute.instanceGroupManagers.create",
    "compute.instanceGroupManagers.delete",
    "compute.instanceGroupManagers.get",
    "compute.instanceGroups.create",
    "compute.instanceGroups.delete",
    "compute.instanceTemplates.create",
    "compute.instanceTemplates.delete",
    "compute.instanceTemplates.get",
    "compute.instanceTemplates.useReadOnly",
    "compute.instances.create",
    "compute.instances.delete",
    "compute.instances.get",
    "compute.instances.setLabels",
    "compute.instances.setMetadata",
    "compute.instances.setServiceAccount",
    "compute.instances.setTags",
    "compute.networks.create",
    "compute.networks.delete",
    "compute.networks.get",
    "compute.networks.updatePolicy",
    "compute.regionOperations.get",
    "compute.subnetworks.create",
    "compute.subnetworks.delete",
    "compute.subnetworks.get",
    "compute.subnetworks.use",
    "compute.subnetworks.useExternalIp",
    "compute.zones.get",
    "iam.roles.create",
    "iam.roles.delete",
    "iam.roles.get",
    "iam.roles.list",
    "iam.roles.update",
    "iam.serviceAccounts.actAs",
    "iam.serviceAccounts.create",
    "iam.serviceAccounts.delete",
    "iam.serviceAccounts.get",
    "iam.serviceAccounts.list",
    "resourcemanager.projects.get",
    "resourcemanager.projects.getIamPolicy",
    "resourcemanager.projects.setIamPolicy",
    "storage.buckets.create",
    "storage.buckets.delete",
    "storage.buckets.get",
    "storage.buckets.getIamPolicy",
    "storage.buckets.setIamPolicy"
  ],
  "stage": "BETA"
}
EOF

gcloud iam roles create GRITProvisioner --project=[projectID] --file=./grit-provisioner-role.json
```

</details>

where `[projectID]` is the ID of your Google Cloud project.

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

<!-- begin: input vars -->
| Name                     | Type                                                     | Required? | Default value | Description                                                                                                                                                                                         |
|--------------------------|----------------------------------------------------------|-----------|---------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `name`                   | `string`                                                 | yes       |               | Name of the deployment. Must be unique in scope of a Google Cloud project.                                                                                                                          |
| `labels`                 | `map(string)`                                            | no        |               | Arbitrary list of `key=value` pairs that are added as labels to resources created by GRIT.                                                                                                          |
| `google_project`         | `string`                                                 | yes       |               | ID of the Google Cloud project.                                                                                                                                                                     |
| `google_region`          | `string`                                                 | yes       |               | Google Cloud region chosen for the deployment.                                                                                                                                                      |
| `google_zone`            | `string`                                                 | yes       |               | Google Cloud zone chosen for the deployment.                                                                                                                                                        |
| `gitlab_url`             | `string`                                                 | yes       |               | URL of GitLab instance.                                                                                                                                                                             |
| `runner_token`           | `string`                                                 | yes       |               | Authentication token of the runner to deploy. See [how to obtain the token](https://docs.gitlab.com/ee/ci/runners/runners_scope.html#create-an-instance-runner-with-a-runner-authentication-token). |
| `runner_machine_type`    | `string`                                                 | no        |               | Machine type for the runner manager instance. If not provided, GRIT uses one of predefined choices based on the value of defined concurrency.                                                       |
| `runner_disk_type`       | `string`                                                 | no        | `pd-standard` | Disk type for the runner manager instance. If not provided, GRIT uses `pd-standard`.                                                                                                                |
| `runner_tags`            | `list(string)`                                           | no        | []            | Tags to register the runner with                                                                                                                                                                    |
| `concurrent`             | `number`                                                 | no        | 50            | Value for `config.toml`'s `concurrent` setting. Defines maximum total number of jobs executed concurrently by the runner.                                                                           |
| `runners_global_section` | `string`                                                 | no        |               | [Allows to customize](#runners_global_section-customization) the global part of `[[runners]]` section in generated `config.toml`.                                                                   |
| `runners_docker_section` | `string`                                                 | no        |               | [Allows to customize](#runners_docker_section-customization) the global part of `[runners.docker]` section in generated `config.toml`.                                                              |
| `capacity_per_instance`  | `number`                                                 | no        | 1             | Maximum number of jobs to be executed concurrently on one autoscaled ephemeral VM.                                                                                                                  |
| `max_instances`          | `number`                                                 | no        | 200           | Maximum number of ephemeral instances (in all possible states) that runner maintains.                                                                                                               |
| `max_use_count`          | `number`                                                 | no        | 1             | Maximum number of jobs executed on a single autoscaled ephemeral VM, after which the VM is marked for deletion.                                                                                     |
| `autoscaling_policies`   | [`list(object)`](#autoscaling_policies-object-structure) | no        |               | List of objects defining autoscaling policies.                                                                                                                                                      |
| `ephemeral_runner`       | [`object`](#ephemeral_runner-object-structure)           | no        |               | Configuration of autoscaled ephemeral VM.                                                                                                                                                           |
<!-- end: input vars -->

### `autoscaling_policies` object structure

| Key                  | Type           | Description                                                                                                                                                                                                                                 |
|----------------------|----------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `periods`            | `list(string)` | List of periods in [unx-cron syntax](https://docs.gitlab.com/runner/configuration/advanced-configuration.html#periods-syntax-1) for which the policy should be applied.                                                                     |
| `timezone`           | `string`       | Identifier of the time zone that should be used to evaluate the periods (for example `UTC`). If not defined, local time zone set for the runner process (so usually the system time zone) is used.                                          |
| `scale_min`          | `number`       | The minimum size of the autoscaling instances fleet. This defines the number of `idle` Taskscaler slots and therefore Fleeting instances that should be sustained for all time to have space for jobs execution. Integer value is expected. |
| `idle_time`          | `string`       | Minimal duration for which the `idle` instances should running, even if not used. Uses Go's time format, for example `1h20m30s`.                                                                                                            |
| `scale_factor`       | `number`       | If used, the number of idle slots are calculated as `scale_factor * in_use_slots`, but not less than defined with `scale_min`. A `float64` value higher than 0 is expected.                                                                 |
| `scale_factor_limit` | `number`       | Usable only when `scale_factor` is in use. If defined, the maximum value of `idle` is calculated with the equation described for `scale_factor`. An integer value is expected.                                                              |

Module expects a list of objects of this type, for example:

```terraform
module "grit-scenario" {
  # (...)
  autoscaling_policies = [
    {
      periods   = ["* * * * sat-sun"]
      scale_min = 1
      idle_time = "30m"
    },
    {
      periods   = ["* * * * *"]
      scale_min = 100
      idle_time = "1h"
    }
  ]
  # (...)
}
```

If not defined by the user, the following default is used:

```terraform
object {
  periods            = ["* * * * *"]
  timezone           = ""
  scale_min          = 10
  idle_time          = "20m0s"
  scale_factor       = 0.2
  scale_factor_limit = 100
}
```

current time, in the defined time zone. If the current time matches any of the periods in the list, the
whole policy entry is used. The last entry that matches current time is taken.

At least one policy with period `"* * * * *` must be defined. This policy is provided as
the scenario's default. To add different settings for the `"* * * * *"` period, define them as
the first entry in the `autoscaling_policies` block. The settings are added after the default,
which means it takes the priority during evaluation.

### `ephemeral_runner` object structure

```terraform
object {
  disk_type    = optional(string, "pd-standard")
  disk_size    = optional(number, 25)
  machine_type = optional(string, "n2d-standard-2")
  source_image = optional(string, "projects/cos-cloud/global/images/family/cos-stable")
}
```

| Key            | Type     | Description                                                                                                     |
|----------------|----------|-----------------------------------------------------------------------------------------------------------------|
| `disk_type`    | `string` | [Disk type](https://cloud.google.com/compute/docs/disks) to be used for autoscaled ephemeral VMs.               |
| `disk_size`    | `number` | Disk size (in GiB) to be used for autoscaled ephemeral VMs. Integer value is expected.                          |
| `machine_type` |          | [Machine type](https://cloud.google.com/compute/docs/machine-resource) to be used for autoscaled ephemeral VMs. |
| `source_image` |          | Source image from which autoscaled ephemeral VMs are started.                                                   |

If not defined by the user, the following default is applied:

```terraform
object {
  disk_type    = "pd-standard"
  disk_size    = 25
  machine_type = "n2d-standard-2"
  source_image = "projects/cos-cloud/global/images/family/cos-stable"
}
```

### `runners_global_section` customization

This setting can be used to add custom configuration to
the [`[[runners]]` section](https://docs.gitlab.com/runner/configuration/advanced-configuration.html#the-runners-section)
of the generated `config.toml` file.

> The `runners_global_section` content is written to the `config.toml`
> file at a location specified in the template. Terraform does not check the syntax.
> If you customize a runner setting that has already been managed by GRIT, GRIT might end with a failing configuration.
> You should deploy the runner before you customize it to review
> the generated `config.toml` file. After you deploy the runner, add the customization.

In the following example:

- The `environment` setting is added to the configuration file.
- The HEREDOC syntax is used to pass multiline content.

```terraform
module "grit-scenario" {
  # (...)
  runners_global_section = <<EOS
  environment = [
    "DOCKER_TLS_CERTDIR=",
  ]
EOS
  # (...)
}
```

### `runners_docker_section` customization

The `runners_docker_section` variable adds customization to the
[`[runners.docker]` section](https://docs.gitlab.com/runner/configuration/advanced-configuration.html#the-runnersdocker-section)
of the generated `config.toml` file.

> The `runners_docker_section` content is written to the `config.toml`
> file at a location specified in the template. Terraform does not check the syntax.
> If you customize a runner setting that has already been managed by GRIT, GRIT might end with a failing configuration.
> You should deploy the runner before you customize it to review
> the generated `config.toml` file. After you deploy the runner, add the customization.

In the following example:

- The `privileged` mode has been enabled for Docker executor.
- Custom volume mapping has been added for every container that starts when jobs are executed.
- The HEREDOC syntax is used to pass multiline content.

```terraform
module "grit-scenario" {
  # (...)
  runners_docker_section = <<EOS
    volumes = [
      "/certs/client"
    ]

    privileged = true
EOS
  # (...)
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
- Initialization of the `google` provider.
- The `runner-manager-external-ip` output.

In both examples, the `runner_token` variable is an external variable marked
as `sensitive`.

The `runner_token` defines the identity of the runner. It provides access to all projects
and jobs for the runner, and must be treated as a secret value.

You should not put the `runner_token` variable into the `*.tf` file
or commit it to any version control system. Instead, it should be passed
to Terraform as a variable, as an environment variable, with the CLI flag or
a variables file. Read more about [how to use input variables with Terraform](https://developer.hashicorp.com/terraform/language/values/variables#assigning-values-to-root-module-variables).

**Simple example**

```terraform
terraform {
  required_version = "~> 1.5"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.12"
    }

    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
}

locals {
  google_project = "your-google-project-ID"
  google_region  = "us-east1"
  google_zone    = "us-east1-b"
}

provider "google" {
  project = local.google_project
  region  = local.google_region
  zone    = local.google_zone
}

variable "runner_token" {
  type      = string
  sensitive = true
}

module "runner-deployment" {
  source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/google/linux/docker-autoscaler-default"

  google_project = local.google_project
  google_region  = local.google_region
  google_zone    = local.google_zone

  name = "gritexample1"

  gitlab_url = "https://gitlab.com"

  runner_token = var.runner_token
}

output "runner-manager-external-ip" {
  value = module.runner-deployment.runner_manager_external_ip
}
```

**Advanced example**

```terraform
terraform {
  required_version = "~> 1.5"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.12"
    }

    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
}

locals {
  google_project = "your-google-project-ID"
  google_region  = "us-east1"
  google_zone    = "us-east1-b"
}

provider "google" {
  project = local.google_project
  region  = local.google_region
  zone    = local.google_zone
}

variable "runner_token" {
  type      = string
  sensitive = true
}

module "runner-deployment" {
  source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/google/linux/docker-autoscaler-default"

  google_project = local.google_project
  google_region  = local.google_region
  google_zone    = local.google_zone

  name = "gritexample1"

  gitlab_url = "https://gitlab.com"

  runner_token = var.runner_token

  concurrent = 1000

  capacity_per_instance = 2
  max_use_count         = 10
  max_instances         = 1000
  autoscaling_policies  = [
    {
      periods = [
        "* * * * sat-sun"
      ]
      scale_min          = 10
      scale_factor       = 0.5
      scale_factor_limit = 500
    },
    {
      periods = [
        "55-59 3 * * sat-sun",
        "00-05 4 * * sat-sun"
      ]
      scale_min          = 200
      scale_factor       = 1.2
      scale_factor_limit = 500
    },
    {
      periods = [
        "* * * * *"
      ]
      scale_min          = 200
      scale_factor       = 1.2
      scale_factor_limit = 500
    }
  ]

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

  ephemeral_runner = {
    disk_type    = "pd-ssd"
    disk_size    = 50
    machine_type = "n2d-standard-8"
    source_image = "projects/cos-cloud/global/images/family/cos-stable"
  }

}

output "runner-manager-external-ip" {
  value = module.runner-deployment.runner_manager_external_ip
}
```

### Terraform execution

Plan and deploy the example as [documented](../../../terraform.md).

### Access to the runner manager instance

After the `terraform apply` step is finished, Terraform prints defined outputs.
In the advanced and simple examples, Terraform prints the external IP of the
instance where runner manager was deployed.

This instance can be now accessed with SSH, using [Google Cloud OS Login](https://cloud.google.com/compute/docs/oslogin)
mechanism.

## Troubleshooting

If things go as expected, you should see a runner-manager VM alongside a few ephemeral VMs in
the Google Cloud project you specified. If that's not the case, try investigating the issue:

```shell
## ssh into the runner-manager VM
## replace the zone, VM and project names as needed
gcloud compute ssh --zone "us-east1-b" "gritexample1-runner-manager" --project "{Your-GCP-Project-ID}"

## Escalate your privileges
$ sudo su
## Check the init script output for any errors
$ less /var/log/cloud-init-output.log
## View the status of the runner service
$ systemctl status gitlab-runner.service -l
```

If you make any changes to your Terraform definitions, you might need to execute the following:

```shell
## Try re-executing key init script
$ /var/lib/cloud/instance/scripts/runcmd
## View the status of the runner service
$ systemctl status gitlab-runner.service -l
```
