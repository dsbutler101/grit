## Developer Certificate of Origin and License

By contributing to GitLab Inc., you accept and agree to the following terms and
conditions for your present and future contributions submitted to GitLab Inc.
Except for the license granted herein to GitLab Inc. and recipients of software
distributed by GitLab Inc., you reserve all right, title, and interest in and to
your Contributions.

All contributions are subject to the
[Developer Certificate of Origin and License](https://docs.gitlab.com/ee/legal/developer_certificate_of_origin).

_This notice should stay as the first item in the CONTRIBUTING.md file._

## Code of conduct

As contributors and maintainers of this project, we pledge to respect all people
who contribute through reporting issues, posting feature requests, updating
documentation, submitting pull requests or patches, and other activities.

We are committed to making participation in this project a harassment-free
experience for everyone, regardless of level of experience, gender, gender
identity and expression, sexual orientation, disability, personal appearance,
body size, race, ethnicity, age, or religion.

Examples of unacceptable behavior by participants include the use of sexual
language or imagery, derogatory comments or personal attacks, trolling, public
or private harassment, insults, or other unprofessional conduct.

Project maintainers have the right and responsibility to remove, edit, or reject
comments, commits, code, wiki edits, issues, and other contributions that are
not aligned to this Code of Conduct. Project maintainers who do not follow the
Code of Conduct may be removed from the project team.

This code of conduct applies both within project spaces and in public spaces
when an individual is representing the project or its community.

Instances of abusive, harassing, or otherwise unacceptable behavior can be
reported by emailing [contact@gitlab.com](mailto:contact@gitlab.com).

This Code of Conduct is adapted from the [Contributor Covenant](https://contributor-covenant.org), version 1.1.0,
available at [https://contributor-covenant.org/version/1/1/0/](https://contributor-covenant.org/version/1/1/0/).

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

#### Go

We lint Go with [golangci-lint](https://golangci-lint.run/).

Run `make lint-go` to run golangci-lint.

We're using a small number of linters as the Go code is minimal and consists only of tests.

#### Documentation

The documentation should follow the [GitLab documentation style guide](https://docs.gitlab.com/ee/development/documentation/styleguide/).

We lint the Markdown files with `vale` and `markdownlint`.

You can read more about configuring each tool in
[GitLab documentation testing](https://docs.gitlab.com/ee/development/documentation/testing/).

Run `make lint-docs` to run both linting tools.
