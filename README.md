# GitLab Runner Infrastructure Toolkit

The GitLab Runner Infrastructure Toolkit (GRIT) is to be the primary entry-point for setting up and managing runner deployments.

## Current state

Experimental. GRIT is in early development and not yet consumable.

Follow [the epic](https://gitlab.com/groups/gitlab-org/ci-cd/runner-tools/-/epics/1) to see progress.

## Working use cases

### Mac Runners on AWS -- development -- fleeting-only

Usage:

```sh
make dev-init
make dev-apply
```

Provide current staging AMI when prompted (e.g. `ami-0fcd5ff1c92b00231`).

Get credentials from `environments/dev/terraform.tfstate` for setting up your runner manager (not automated yet).

## Intended use cases

Setting up and managing...
1. a production runner deployment with autoscaling.
2. autoscaling on Google, AWS and Azure.
3. a runner with blue-green deployment.
4. Linux, MacOS and Windows runners.
5. instance, group and project runners.
6. a Fleeting stack for use while developing runner.
7. ephemeral stacks for end-to-end testing.
