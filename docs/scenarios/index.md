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

## Prerequisites

### Terraform

GRIT is a library of Terraform modules. To use it, you must have a working Terraform setup.

Depending on your usage, execution of Terraform might typically be done with automation.
For basic experimentation on local machines, an installation of Terraform CLI is
required.

Read more about [how to install Terraform CLI](https://developer.hashicorp.com/terraform/install).

## Available scenarios

### Google

Setups deployed to Google Cloud

#### Google Cloud integration prerequisites

##### Google Cloud SDK

Terraform, no matter what execution method was chosen, requires support for
Google Cloud SDK. And for that a setup of credentials for authenticating
requests to Google Cloud API is required.

Details of how to set this up depend on the chosen method. The most simple
one using local Terraform CLI execution works best if `gcloud` command is
installed locally as well.

To use the Google Cloud scenario templates, you must have:

- [Google Cloud CLI installed](https://cloud.google.com/sdk/docs/install).
- Credentials to authenticate to the Google Cloud API. For more information, see [Initializing the gcloud CLI](https://cloud.google.com/sdk/docs/initializing).

##### Google Cloud project

For the scenarios from the Google group, access to a Google Cloud project is required.

You should use a dedicated Google Cloud project for CI/CD workloads. This is not a strict
requirement, but it provides the following advantages:

- Separation of context from other workloads.
- Separation of access and permissions, which limits how resources are accessed.
- Better observability.
- Easier cost analysis.

If you have configured Workload Identity Federation, you should use a different Google Cloud project
than the one where you created the Workload Identity Pool and provider.

##### Billing account linked with Google Cloud project

To enable API services required by GRIT, a billing account in Google Cloud must be linked
with the project chosen for GRIT deployments.

For more information, see [Check if billing is enabled on a project](https://cloud.google.com/billing/docs/how-to/verify-billing-enabled#confirm_billing_is_enabled_on_a_project).

#### Linux

Scenarios to deploy GitLab Runner configurations that execute CI/CD jobs on Linux platform.

1. [Docker Autoscaler default](google/linux/docker-autoscaler-default)
