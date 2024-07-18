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
