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
  source = "../../grit/modules/test"

  manager_provider  = "ec2"
  capacity_provider = "none"

  gitlab_project_id         = "YOUR_PROJECT_ID
  gitlab_runner_description = "grit-runner"
  gitlab_runner_tags        = []
}
```

#### GKE Kubernetes Runner

```terraform
module "gke-kubernetes-runner" {
  source = "../../grit/modules/test"

  manager_provider  = "helm"
  capacity_provider = "gke"

  gitlab_project_id         = "YOUR_PROJECT_ID"
  gitlab_runner_description = "grit-runner"
  gitlab_runner_tags        = []
}
```

### Dev

Dev use cases setup a piece of the runner infrastructure which is
convienient for development and debugging. For example setting up an
Instance Group for Fleeting on the system of choice, outputting just
the credential necessary to access the raw resources.

#### Mac Runners on AWS -- Fleeting Instance Group

```terraform
module "my-experimental-mac-machines" {
  source = "../grit/modules/dev"

  fleeting_provider = "ec2"
  os                = "macos"
}
```

Get credentials from `environments/dev/terraform.tfstate` for setting up your runner manager (not automated yet).

## Value of GRIT

The infrastructure-as-a-library approach of GRIT provides value across
many personas and use-cases.

### Testing

- Easier to setup a demo with non-trivial runner infrastructure.
- A single entry-point for discovering and learning about runner configuration.
- A common test library for verifying changes to runner don't break user infrastructure.

### Development

- Quickly setup the parts of the stack that a developer *isn't* working on, so they can focus on the part they *are*.
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
