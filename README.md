Demo: https://youtu.be/sWugZ_eW5nQ

# GitLab Runner Infrastructure Toolkit (GRIT)

The GitLab Runner Infrastructure Toolkit (GRIT) is a library of
Terraform templates for deploying GitLab Runner and managing its
lifecycle. It covers everything from a single runner deployment to
complex autoscaling configurations. It embodies the best-practices
for configuration and operation of runner.

## Current state

Alpha. There are a few things that work. But they are likely to change
quite a bit.

Follow [the epic](https://gitlab.com/groups/gitlab-org/ci-cd/runner-tools/-/epics/1) to see progress.

## Working use cases

### Test

Test use cases setup a working runner stack and register it to a
GitLab instance. The result is a working system, but one which
uses convenient defaults and is not necessarily production grade.

#### Single EC2 Shell Runner

```terraform
module "single-ec2-shell-runner" {
  source = "modules/aws/test"

  manager_service  = "ec2"
  fleeting_service = "none"

  gitlab_project_id         = "YOUR_PROJECT_ID"
  gitlab_runner_description = "grit-runner"
  gitlab_runner_tags        = []
  name                      = "test-name"
}
```

#### GKE Kubernetes Runner

```terraform
module "gke-kubernetes-runner" {
  source = "modules/aws/test"

  manager_service  = "helm"
  fleeting_service = "gke"

  gitlab_project_id         = "YOUR_PROJECT_ID"
  gitlab_runner_description = "grit-runner"
  gitlab_runner_tags        = []
  name                      = "test-name"
}
```

### Dev

Dev use cases setup a piece of the runner infrastructure which is
convenient for development and debugging. For example setting up an
Instance Group for Fleeting on the system of choice, outputting just
the credential necessary to access the raw resources.

#### Mac Runners on AWS -- Fleeting Instance Group

```terraform
module "my-experimental-mac-machines" {
  source = "modules/aws/dev"

  fleeting_service = "ec2"
  fleeting_os      = "macos"
  ami              = "ami-12345"
  instance_type    = "mac2.metal"
}
```

Get credentials from `environments/dev/terraform.tfstate` for setting up your runner manager (not automated yet).

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

The GRIT API is defined by `variables.tf` in the directories
`modules/dev`, `modules/test` and `modules/prod`. These variables are
the parameters for configuring each stage of GRIT.

The `dev` stage should be considered volatile. It's a set of
convenience configuration for developing and debugging runner. They
should work but might change as needs of developers change. These
templates do not produce a full working runner deployment.

The `test` stage will be relatively stable. It's all the ways that
GRIT can deploy runners but with lots of convenient defaults so they
can be setup will few required parameters. These templates register
runner to a GitLab instance so they will produce a working runner
deployment.

The `prod` stage templates will each be associated with a maturity
designation of `alpha`, `beta` or `stable`. The goal of all `prod`
templates is to be `stable`.

The `stable` templates are the set of battle-hardened templates that
GRIT authors have experience running. They will be backward
compatible, meaning they will not be refactored in a way that causes
unnecessary resource destruction. Especially since that would likely
cause an outage! When backward incompatable changes are necessary
advance warning will be given in manner compatible with existing
GitLab deprecation policies.

## Contributing

GRIT is currently early in the development stage so if you want to
contribute, reach out to us through Slack or open an Issue. There's
lots to do, but you might need a little help getting started.

### General Guidance

The GRIT codebase should conform to [Google's best practices for using
Terraform](https://cloud.google.com/docs/terraform/best-practices-for-terraform).

The goal of GRIT is decompose runner infrastructure sufficiently that
there is little to no repetition. The `environments` folders are
examples of how to use GRIT to setup a runner deployment. The
`modules/dev` and subfolders are all the defaults that are applied to
make `dev` templates convenient to
use. E.g. `modules/dev/ec2/macos/macos.tf` contains a reference to a
default AMI. Same goes for `modules/test` and `modules/prod`. The
`modules/internal` directory structure contains the implementation
details of each template, separate by cloud provider and then
operating system. These should not be used outside the GRIT modules.

### Outputs

Because GRIT handles a widely expanding combination of configurations
outputs are exported from each layer as a map. This avoids the need to
repeat each output at every layer.

For example EC2 MacOS outputs are surfaced like this:

```terraform
output "output_map" {
  description = "Outputs from the Fleeting Instance Group"
  value = tomap({
    "ssh_key_pem"                                = module.instance_group.ssh_key_pem,
    "fleeting_service_account_access_key_id"     = module.instance_group.fleeting_service_account_access_key_id,
    "fleeting_service_account_secret_access_key" = module.instance_group.fleeting_service_account_secret_access_key,
  })
}
```

And bubbled up like this:

```terraform
output "output_map" {
  description = "Outputs from EC2 resources"
  value = tomap({
    "macos" = try(module.macos[0].output_map, null),
  })
}
```

For resources which are not created based on `counts` fields
the `try` block will emit their output map as `null`.

This allows users to access the ssh key via references like this:

```terraform
output "my-runner-ssh-key" {
  value = module.my-runner.ec2.macos.ssh_key_pem
}
```

And prevents outputs referencing resources not created. E.g. AKS
clusters when deploying to GKE.
