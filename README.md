# GitLab Runner Infrastructure Toolkit

The GitLab Runner Infrastructure Toolkit (GRIT) is to be the primary entry-point for setting up and managing runner deployments.

## Current state

Experimental. GRIT is in early development and not yet consumable.

Epic: https://gitlab.com/groups/gitlab-org/ci-cd/runner-tools/-/epics/1

## Intended use cases

Setting up and managing...
1. a production runner deployment with autoscaling.
2. autoscaling on Google, AWS and Azure.
3. a runner with blue-green deployment.
4. Linux, MacOS and Windows runners.
5. instance, group and project runners.
6. a Fleeting stack for use while developing runner.
7. ephemeral stacks for end-to-end testing.
