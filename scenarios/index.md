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

### Google cloud

1. [Linux - Docker Autoscaler default](google/linux/docker-autoscaler-default)

### AWS

1. [Linux - Docker Autoscaler default](aws/linux/docker-autoscaler-default)
