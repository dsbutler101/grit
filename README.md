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

Download the [latest release](https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit/-/releases)
and reference the test or prod modules in your Terraform configuration.

### Test and Prod Modules

GRIT provides two types of modules, test and prod. Test modules
provide lots of defaults and have no support promises or
restrictions. They are useful for local development and automated
testing.

Prod modules have associated support levels and promises around
backward compatibility. Support levels follow the
[GitLab definition](https://docs.gitlab.com/ee/policy/experiment-beta-support.html)
of experimental, beta and GA. See APIs and Guarantees below for more
information.

Consumable modules are organized as follows:

1. `modules` directory
1. provider (e.g. `aws`, `google`, `azure`)
1. module (e.g. `vpc`, `iam`, `fleeting`, `runner`)
1. type (e.g. `prod` or `test`)

Example module source path:

```hcl
module "my-runner" {
  source = "grit/modules/aws/runner/prod"
  ...
}
```

All consumable module paths end in either "test" or "prod". Do not
directly reference any modules within "internal" folders.

### Composable Modules

The primary module is `runner`, which can be used by itself
([example main.tf](examples/test-shell-runner-only-ec2/main.tf)).
Required and optional inputs are documented in the `variables.tf` file
([example prod variables.tf](modules/aws/runner/prod/variables.tf)).
Outputs are documented in the `outputs.tf` file
([example prod outputs.tf](modules/aws/runner/prod/outputs.tf)).

Optional modules are available to set up additional configuration for
runner which can be fed into the `runner` module. For example, the
`gitlab`, `vpc`, `iam`, `fleeting` and `runner` modules can create a
fully autoscaling runner on its own VPC and automatically register it
to a GitLab project ([example main.tf](examples/prod-docker-autoscaler-ec2/main.tf)).
The outputs of each optional module are exactly what is required as input to the
`runner` module, so they should fit together easily.

### Examples

- [Test shell runner on EC2](examples/test-shell-runner-only-ec2/main.tf)
- [Production `docker-autoscaler` configuration on EC2](examples/prod-docker-autoscaler-ec2/main.tf)

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
- Easier to discuss and diagnose configuration issues. Users can say "I'm on GRIT 16.5 and I'm seeing this ..." and support will know exactly how they are setup.

## API and Guarantees

The GRIT API is defined by `variables.tf` and `output.tf` in the
`test` and `prod` directories.

The `test` type modules provide all the ways that GRIT can deploy
runners with lots of convenient defaults so they can be setup will few
required parameters.

The `prod` type modules will each be associated with a support
designation of `experimental`, `beta` or `ga`. The end-goal is to have
all `prod` modules be `ga`. In general, experimental modules are new
or just used in tests and development. Beta modules are at least
dogfooded by GitLab internally. And GA modules are used by GitLab
customers. Some modules like `runner` provide support levels on a
per use case basis.

The `ga` modules are the set of battle-hardened modules that GRIT
authors have experience running. They will be backward compatible,
meaning they will not be refactored in a way that causes unnecessary
resource destruction. Especially since that would likely cause an
outage! When backward incompatible changes are necessary advance
warning will be given in manner compatible with existing GitLab
deprecation policies.

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

GRIT is currently early in the development stage so if you want to
contribute, reach out to us through Slack or open an
[Issue](https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit/-/issues). There's
lots to do, but you might need a little help getting started.
