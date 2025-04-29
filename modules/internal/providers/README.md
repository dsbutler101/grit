# Providers meta module

This module has a reference to every unique provider used anywhere in GRIT.
It is not to be referenced from any other module.

The purpose of this is to provide one way to download all dependencies
in a single operation rather than downloading them for every module.
Then we can use either the plugin cache or `--plugin-dir` to ensure
all dependencies are available before each `terraform init` runs.
See also [provider installation docs](https://developer.hashicorp.com/terraform/cli/config/config-file#provider-installation).

Each provider plugin can be quite a large binary, potentially many hundreds
of MB. There is no way to use the terraform plugin cache concurrently either,
so we have to ensure we download all dependencies once and reuse them.

See also [terraform parallel init issue](https://github.com/hashicorp/terraform/issues/25849)
which describes the issues with init concurrency.
