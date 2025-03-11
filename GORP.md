# GRIT Oversight and Recommendation Panel (GORP)

The objective of GRIT is to provide the widest possible collaboration
space in the domain of operating GitLab Runner on public clouds. As
such the configuration space is enormous: the dot product of all
clouds, all compute services, all GitLab Runner executors and all
Fleeting plugins. The GitLab Runner team cannot build or maintain all
configurations. And configurations are best maintained by users that
are dogfooding them. So GRIT configurations are implemented lazily by
the users that use them.

This model risks diverging styles and feature parity between the
different cloud providers and use cases, which makes GRIT more
difficult to understand and maintain overall. This is a real problem
observed during the first year of GRIT. Some provider and service
details differ between modules, but overall they must all be
coherent. To provide the strong central guidance necessary to achieve
coherence, maintain simplicity and consistent style we have
established the GRIT Oversight and Recommendation Panel (GORP).

## Shared Responsibility

GRIT uses a model of shared responsibility between Authors and
Operators.

Authors are responsible for building and maintaining the GRIT
terraform modules, tools, tests and their release. They must follow
the guidelines set out by the GORP. Authors must identify the
[support stage](https://docs.gitlab.com/policy/development_stages_support/)
for each configuration they provide (`experimental`, `beta` or
`ga`). Configurations which have not been evaluated for support level
will be marked as `none` but may technically work.

Operators are responsible for the runners they deploy. They must
understand the implications of the GRIT configurations they use. They
must explicitly declare the support stage they require in order to
ensure they are getting what they expect. This is done via the
`min_support` variable. Operators must only depend on released
versions of GRIT unless prior arrangements are made.

The GORP is responsible for providing Authors with clear, prescriptive
guidance for how to make changes to GRIT which fit with the expected
style and architecture. GORP review is NOT required for all
changes. GRIT changes are proposed by a Contributors following the
GitLab DRI process which ensures individuals are both responsible
and empowered. Merge requests are approved by GRIT repository
approvers in the usual way. Contributors must follow the GORP
guidelines but they are free to interpret them and should not be
blocked by GORP guidelines. Only changes to the GORP responsibilities
and guidelines (this document) must be reviewed by a GORP member.

## GORP Guidelines

These guidelines are provided by the GRIT Oversight and Recommendation
Panel (GORP). Please follow these guidelines when making or approving
changes to GRIT. If these guidelines need clarification, additions or
updating, please open an issue and assign it to the GORP.

The GORP is currently composed of @josephburnett, @amknight and
@tmaczukin.

### Composible Architecture

GRIT provides a vertically layered set of tools for operating GitLab
Runner. The lowest level are the individual Terraform modules found
within the `/modules` folder. These modules are also horizontally
composible, each encapsulating some aspect of the runner environment
for a given cloud provider. For example `fleeting` capacity is
provided as a separate module from `runner`. The runner will use
Fleeting or not based on the executor configured. The outputs of each
module are designed to compose seamlessly into the inputs of other
modules which depend on them.

One layer higher are scenarios found in the `/scenarios` folder.
These show how the modules are combined into common and complete
configurations. These are usually composed around a specific provider,
service and executor.

Higher than scenarios are tools for deploying and maintaining runners
throughout their lifecycle. This is the Deployer binary found in the
`/deployer` folder. Deployer is used to stand up a runner
configuration, determine its status and send it lifecycle
signals. This is done by means of the GitLab Runner `wrapper` command
and gRPC service.

Above Deployer are the CI Steps which can be used to run Deployer in
GitLab CI. CI Components use the steps to produce a full runner
deployment pipeline. The pipeline can either target the same GitLab
instance (self-installing) or a separate GitLab instance (e.g. from an
"ops" instance). There are not yet CI Steps or Components in GRIT.

Similar to scenarios are reference architectures. These are also
complete configurations but are fewer and more carefully
currated. Each reference architecture also includes guidance for
scalability and specific characteristics which make it desirable, as
well as its limitations. Reference architectures may include Deployer
to setup long-lived blue-green deployments. There are not yet
reference architectures in GRIT.

#### Perfect Fit

The outputs of one module must fit perfectly into the inputs of
modules which depend on it. A perfect fit means the named outputs and
types of a module (such as `vpc`) are exactly the fields and types
required by other modules (such as `runner` and `fleeting`).

For example if `vpc` provides these outputs:

```hcl
output "id" {
  type = string
  ...
}

output "subnet_ids" {
  type = list(string)
  ...
}
```

Then both `runner` and `fleeting` which depend on `vpc` will accept a
structure named `vpc` with identical fields:

```hcl
variable "vpc" {
  type = object({
    id         = string
    subnet_ids = list(string)
  })
  ...
}
```

This ensures that top-level modules remain composible. The outputs of
a dependency can be stored in a `local` structure which can be passed
as-is to all dependants. It requires a little extra work to update all
GRIT dependants inputs when updating outputs.

#### Thin Scenarios

Scenarios must be composed primary of lower-level modules. They must
be as thin as possible. This allows Operators to begin with a scenario
and evolve their configuration by "exploding" the scenario into
individual modules without signficantly changing their setup. If some
domain-specific logic is worth adding to a scenario, it's probably
worth adding to the low-level module. For example automatic selection
of instance size for runner managers should happen in the `runner`
module, not in the `google-docker-autoscaler` scenario.

#### Configuration Agnostic Deployer

Deployer must not be aware of the details of the runner configuration
it is controlling. It accepts only the folder which contains the
configuration and reads the output address of the runner
managers. This allows Deployer to be reused in several deployment
styles and systems, such as blue-green or rolling deployments. It also
allows rollout of large runner configuration changes without changing
the high-level deployment infrastructure.

#### Dogfood Steps and Components

GRIT CI Steps and Components should be the same ones we use for
GitLab.com Hosted Runners. Dogfooding is the best way to support
self-hosted users of GRIT.

#### Consistent and Predictable Structure

The GRIT folder structure must be very consistent and predictable so
the codebase is easy to navigate. All consumable modules must be in
the `modules` directory. The next directory layer must be the
provider. We separate modules by provider so consumers are not forced
to configure any providers they are not using.

The next layer must be a series of logical modules. The primary module
must be `runner`. The rest can be a decomposition of runner setups on
the provider platform (for example, the `aws` folder contains `iam`,
`vpc`, `fleeting`, and `runner` modules).

#### Top-Level Modules

Top-level modules in a provider should represent highly-decoupled or
optional configuration aspects of runner. For example, `fleeting` and
`runner` are coupled only by access credential and the name of the
instance group, so they are separate modules. VPC details are optional
because some users bring their own, so `vpc` is a separate
module. Users that bring their own VPC need only create a matching
input structure to plug into the rest of the GRIT modules.

It is preferable to have a similar decomposition of top-level modules
across providers. For example `vpc` should have the same name within
`aws` and `google`. Additionally the inputs and outputs should be
named the same whenever possible.

### General Style

In general GRIT Terraform modules should be simple and unsurprising.

#### Best Practices

In the absense of other guidelines, the GRIT codebase must conform to
[Google's best practices for using Terraform](https://cloud.google.com/docs/terraform/best-practices-for-terraform).

#### Location Agnostic

GRIT modules should not have a default region or zone. This prevents
accidentally omission creating resources in the wrong location. And It
prevents hot-spotting everything the default localtion.
