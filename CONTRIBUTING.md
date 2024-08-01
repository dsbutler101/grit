## Contribute to GRIT

### General Guidance

The GRIT codebase must conform to
[Google's best practices for using Terraform](https://cloud.google.com/docs/terraform/best-practices-for-terraform).

The goal of GRIT is to decompose runner infrastructure sufficiently that
there is little repetition. GRIT uses composable modules to reduce the
complexity of each module. See the
[Zen of Fabric](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/blob/master/CONTRIBUTING.md#the-zen-of-fabric).

The structure should be very consistent and predictable so the
codebase is easy to navigate. All consumable modules must be in the
`modules` directory. The next directory layer must be the provider.
We separate modules by provider so consumers are not forced to configure
any providers they are not using.

The next layer must be a series of logical modules. The primary module
must be `runner`. The rest can be a decomposition of runner setups on
that provider platform (for example, the `aws` container `iam`, `vpc`, `fleeting`,
and `runner` modules). Consider creating a separate module for aspects
with low coupling where consumers might want to bring their own.

### Tests

Testing Terraform requires access to cloud providers.

We have very few unit tests here because we must be able to plan and validate the
Terraform output against real resources.

Each set of tests is run in Go test files alongside the Terraform code.
We use [Terratest](https://terratest.gruntwork.io/) to interact with Terraform.

#### Integration Tests

Terraform requires installation of AWS and GCP command line interfaces.

For AWS you can follow the [Terraform AWS pre-requisites](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/aws-build#prerequisites).

For GCP you can follow the [Terraform GCP pre-requisites](https://developer.hashicorp.com/terraform/tutorials/gcp-get-started/google-cloud-platform-build#prerequisites).

As a minimum configure the following environment variables to run the tests:

```shell
export AWS_ACCESS_KEY_ID=<your-aws-access-key>
export AWS_SECRET_ACCESS_KEY=<your-aws-secret-key>
export AWS_REGION=<your-aws-region>
export GOOGLE_PROJECT=<your-google-project-id>
export GOOGLE_REGION=<your-google-region>
export GOOGLE_ZONE=<your-google-zone>
```

For GCP this assumes [application default credentials](https://cloud.google.com/docs/authentication/application-default-credentials).

Run the tests with `make test`.

#### End-To-End Tests

End-to-end tests use the modules to create and register runner managers against the
[GRIT End-to-End Test Project](https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit-e2e).
A pipeline is created for the runner and each job is required to complete successfully.

The tests use a fixed GitLab project ID so they can only be run locally by team members.

E2E test code is in the `e2e` directory.

In addition to the above integration test setup, the following environment variables must be set:

```shell
export CI_JOB_ID=<numeric-job-id> # used for naming the resources
export GITLAB_TOKEN=<gitlab-pat-with-api-scope> # used for GitLab automation
export RUNNER_TOKEN=<gitlab-runner-token> # runner token for GRIT e2e project
```

Run the tests with `make e2e-test`.

### Linting

#### Terraform

We check the Terraform style with `terraform fmt` and validate `variables.tf`
files with a `go` test to ensure variable definitions are in the correct place.

Run `make terraform-fmt-check` to check Terraform formatting and
`make lint-terraform` to lint variables.

#### Documentation

The documentation should follow the [GitLab documentation style guide](https://docs.gitlab.com/ee/development/documentation/styleguide/).

We lint the Markdown files with `vale` and `markdownlint`.

You can read more about configuring each tool in
[GitLab documentation testing](https://docs.gitlab.com/ee/development/documentation/testing/).

Run `make lint-docs` to run both linting tools.
