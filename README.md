# GitLab Runner Infrastructure Toolkit (GRIT)

The GitLab Runner Infrastructure Toolkit (GRIT) is a library of
Terraform modules for deploying GitLab Runner and managing its
lifecycle. It covers everything from a single runner deployment to
complex autoscaling configurations. It embodies the best-practices
for configuration and operation of runner.

[Watch this demo](https://youtu.be/sWugZ_eW5nQ) for more details.

## Current State

Experimental. There are a few consumers and the first production use
case will be in beta soon.

Follow [the epic](https://gitlab.com/groups/gitlab-org/ci-cd/runner-tools/-/epics/1) to see progress.

## Usage

Download the [latest release](https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit/-/releases) and reference the appropriate modules in your Terraform configuration.

### Module Structure

GRIT provides modules organized by cloud provider and functionality. Each module is designed to be composable, allowing you to tailor configurations to your specific needs.

Modules have associated support levels and promises around
backward compatibility. Support levels follow the
[GitLab definition](https://docs.gitlab.com/ee/policy/experiment-beta-support.html)
of `experimental`, `beta` and `GA`. See APIs and Guarantees below for more
information.

Modules are organized as follows:

1. `modules` directory
1. Provider (for example, `aws`, `google`, `azure`)
1. Module (for example, `vpc`, `iam`, `fleeting`, `runner`)

Example module source path:

```hcl
module "my-runner" {
  source = "grit/modules/aws/runner"
  ...
}
```

### Composable Modules

The primary module is `runner`, which can be used by itself
([example main.tf](examples/test-shell-runner-only-ec2/main.tf)).
Required and optional inputs are documented in the `variables.tf` file
([example prod variables.tf](modules/aws/runner/variables.tf)).
Outputs are documented in the `outputs.tf` file
([example prod outputs.tf](modules/aws/runner/outputs.tf)).

Optional modules are available to set up additional configuration for
runner which can be fed into the `runner` module. For example, the
`gitlab`, `vpc`, `iam`, `fleeting` and `runner` modules can create a
fully autoscaling runner on its own VPC and automatically register it
to a GitLab project ([example main.tf](examples/docker-autoscaler-ec2-deployed-with-gitlab-ci/main.tf)).
The outputs of each optional module are exactly what is required as input to the
`runner` module, so they should fit together easily.

### Examples

- [Shell runner on EC2](examples/test-shell-runner-only-ec2/main.tf)
- [Docker Autoscaler configuration on EC2](examples/docker-autoscaler-ec2-deployed-with-gitlab-ci/main.tf)
- [Operator on GKE](examples/test-runner-gke-google/main.tf)

## Value of GRIT

The infrastructure-as-a-library approach of GRIT provides value across
many personas and use-cases.

### Testing

- Easier to set up a demo with non-trivial runner infrastructure.
- A single entry-point for discovering and learning about runner configuration.
- A common test library for verifying changes to runner don't break user infrastructure.

### Development

- Quickly set up the parts of the stack that a developer *isn't* working on, so they can focus on the part they *are*.
- Easily reproduce production issues by setting up an identical stack.

### Production

- Dogfooding GRIT makes GitLab SaaS Runners more transparent. Users can see exactly how their jobs are being handled.
- Users can more easily contribute back to SaaS Runners because they can test changes to infrastructure before submitting a merge request.
- GRIT will power Dedicated Runners, allowing the product to share experience directly with the SaaS Runner team via library.
- GitLab Cells will benefit from the same sharing of experience and small, orthogonal configuration surface.

### Self-Hosted

- Users managing their own runner stacks get access to best-practices by default.
- Through a shared library experience we can learn from our self-hosted customers who develop unique experience and expertise.

### Wide Reach

- With a regular cadence of release and upgrade, best practices and bug fixes can be adopted widely.
- Easier to discuss and diagnose configuration issues. Users can say "I'm on GRIT 16.5 and I'm seeing this ..." and support will know exactly how they are set up.

## API and Guarantees

The GRIT API is defined by `variables.tf` and `outputs.tf` in each module directory.

Modules are associated with a support designation of `none`, `experimental`, `beta`, or `GA`. The goal is for all 
modules to reach the `GA` status.

- **None**: Modules with no support guarantees, primarily for testing and development.
- **Experimental**: New modules or those used mainly in tests and development.
- **Beta**: Modules that are at least dogfooded by GitLab internally.
- **GA (generally available)**: Modules used by GitLab customers. Maintain backward compatibility to prevent unnecessary resource destruction. Any necessary backward-incompatible changes will be communicated in advance, adhering to existing GitLab deprecation policies.

### Supported AMIs and AMI deprecation

The `modules/aws/ami_lookup` module provides AMIs for GRIT configurations.
AMIs remain publicly available for up to two years if they are either:

- Listed in `modules/aws/ami_lookup/manifest.json` in a tagged release.
- Created after the latest GRIT release.

Public AMIs are removed if they are either created:

- Before the latest GRIT release and have no association with a tagged GRIT version.
- More than two years ago.

If your GRIT version uses AMIs that are removed, upgrade to the latest version.

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

GRIT is currently early in the development stage so if you want to
contribute, reach out to us through Slack or open an
[Issue](https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit/-/issues). There's
lots to do, but you might need a little help getting started.
