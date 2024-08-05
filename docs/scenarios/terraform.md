---
stage: Verify
group: Runner
info: To determine the technical writer assigned to the Stage/Group associated with this page, see https://handbook.gitlab.com/handbook/product/ux/technical-writing/#assignments
---

# Terraform execution

When the Terraform code is ready, in the directory where the `*.tf` files are created,
complete the following steps to execute the Terraform code:

## Initialize

Use the `init` call to initialize the directory that contains the Terraform configuration files.
The `init` call downloads all providers and external modules referenced by the code.:

```shell
terraform init
```

You must run this command when:

- You first use Terraform in this directory.
- Version definitions in `terraform` block are changed.
- The GRIT code has been updated.

## Plan

Use the `plan` call to compare the local Terraform state stored in the file with the
resources already deployed (e.g. in Google Cloud).

```shell
terraform plan -var runner_token="your-glrt-runner-token" -out plan.out
```

The `plan` call prints all changes that will be provisioned and stores them
in the `plan.out` file that you use in the last step to execute Terraform.

In the above we've explicitly passed in the variable `runner_token` with the
value `your-glrt-runner-token`, by using the `-var` flag for `terraform`. The
`-var` flag can be repeated multiple times with different variable names and
values.

Alternatively, you can export a `TF_VAR_runner_token="glrt-runner-token-here"`
environment variable before you run the `plan` call, for example:

```shell
export TF_VAR_runner_token="glrt-runner-token-here"
terraform plan -out plan.out
```

## Apply the plan

Use the `apply` call to execute the provisioning steps in the `plan.out` file.
Provisioning is executed with calls to the underlying infrastructure via the
used providers (e.g. to Google Cloud API, or to a Kubernetes API).

```shell
terraform apply plan.out
```
