---
status: proposed
creation-date: "2024-07-25"
authors: [ "@tmaczukin", "@josephburnett", "@amknight" ]
---

# GRIT zero-downtime deployments strategy

[[_TOC_]]

## Summary

[GRIT](https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit/) is a library of templates,
intended to provide a comprehensive set of tools for deploying and managing
[GitLab Runner](https://gitlab.com/gitlab-org/gitlab-runner) in different scenarios across
different environments.

Today it contains [low-level terraform modules](../../../modules), providing building blocks
for what could be used to deploy the runner, and a [higher-level terraform modules](../../../scenarios)
providing tested configurations built from these building blocks.

What GRIT is still missing is a solution for managing Deployments in scale.

This document drafts the design proposal for how we can do it, following GRIT's modular
philosophy of "use only what you need".

## Glossary

Let's first define a common wording to make the design description easier to understand and to
apply to a specific use case.

**GitLab Runner** - is [the product](https://gitlab.com/gitlab-org/gitlab-runner), part of the GitLab CI/CD
system that executes jobs supporting different execution strategies and reports the results back to GitLab.

**Runner Process** - is the system process running the GitLab Runner binary.

This is what asks GitLab for jobs, prepares and schedules their execution and reports the result back
to GitLab using the dedicated API.

**Runner Process Wrapper** - is a dedicated GitLab Runner command that allows to run the main Runner Process
in a managed way, with gRPC server for external integrations.

**Graceful Shutdown** - a GitLab Runner builtin mechanism to stop Runner Process without disrupting
executed jobs.

When it's initiated, Runner Process continues to work, but stops asking for new jobs. Execution of jobs
already handled by Runner Process is continued until all of them are finished. Once the last job is
removed from Runner Process' memory, Runner Process exits.

Graceful Shutdown may be executed on the Runner Process using the `SIGQUIT` signal on the Unix like
platforms. It can be also initiated by a dedicated gRPC call if Runner Process is started and managed
through the Runner Process Wrapper.

**Forceful Shutdown** - a GitLab Runner builtin mechanism to stop Runner Process immediately.

That interrupts and cancels jobs being executed. As it creates an unpleasant user experience,
it should be used only in specific situations, when really necessary.

The Runner Process Wrapper exposes a dedicated gRPC method to initiate the Forceful Shutdown.

**Runner Manager** - is the instance where the Process is executed.

A Runner Manager contains only one Runner Process with one `config.toml` file.

**Shard** - is a specific, dedicated offering with its specific configuration.

Shard is an abstract concept, and it may represent, for example:

- For Hosted Runners for GitLab.com, a specific offering like `saas-linux-small-amd64` providing
  Linux runners on x86_64 architecture with the "small" instance sizes that execute jobs.

- For Dedicated Runners that would be a specific customer installation.

- For self-managed GitLab Runner setups, it may be a Runner provided for a specific team in the organization.

A Shard can use multiple Runner Managers, all using exactly the same configuration, to provide
High Availability and horizontal scaling of the capacity.

A Shard can have several Deployment Versions.

**Deployment Version** - is the specific version of the shard configuration.

Deployment Versions provide the same shard offering, but may vary with the details of the configuration
(for example, updating the capacity of the Runner Manager provided by the `concurrent` setting), version
of the installed GitLab Runner binary, definition of the Instance Group for fleeting autoscaling etc.

Excluding the time of performing the Deployment operation, only one Deployment Version should be
active and executing jobs.

**Deployment** - an operation during which we transition a Shard from one Deployment Version to another
Deployment Version.

This is the concept that this design document focuses on.

**Zero-downtime Deployment** - a Deployment strategy that makes the operation totally invisible for
the users and that doesn't affect or disrupt jobs.

**Deployer** - a dedicated binary provided by GRIT project.

This binary will be responsible for executing specific common steps of the Deployment process.

## Goals and requirements

1. The design should be modular, just like the rest of GRIT project.

   Except for what needs to be updated in the products themselves (we will cover a necessary update
   in GitLab Runner, for example), all Deployment mechanics should be put into low and high-level modules.

   This will allow users to use the mechanism in the same way how GRIT is used today - chose only
   the modules that are needed in the specific scenario.

1. Except for changes in the products, all modules of the Deployment mechanism should be hosted
   within the GRIT repository.

1. When possible, we should take usage of GitLab CI/CD tools - steps and components - when developing the
   modules.

1. The designed mechanism must provide support for Zero-downtime Deployment strategy.

1. The designed mechanism must provide support for Prometheus based monitoring.

1. Deployment mechanism should not enforce any specific layout of the Terraform project.

   Users should be free to decide how they want to structure their Terraform code, and the
   deployment mechanism should have as few requirements as possible.

## Design

GRIT's Deployment mechanism will be built upon two things:

1. Teffaform convention on how to manage GitLab Runner Deployments.
1. A Deployer binary to support the desired Deployment process.

All Shard configuration should be provided through Terraform code. Form of the project, granularity
of the used GRIT modules is totally arbitrary and users have a free will on how to design it.

Because the Runner Process is a long-running process with specific termination requirements (described
in more details [below](#runner-process-wrapper)), to achieve a true Zero-downtime Deployment
Strategy, we must ensure that each Deployment Version is handled independently.

We will not enforce any specific project layout. However, we expect that when performing the Deployment,
Terraform configuration of each Deployment Version will be stored in a separate, dedicated directory
and will use a separate, dedicated Terraform State.

The Deployment Version in use when approaching the Deployment operation will be called "old".

The Deployment Version that Deployment operation wants to enroll will be called "new".

The environment in which Deployer commands are executed must have:

- `terraform` command available in system's `PATH`. If this command is not discoverable through `PATH`
  or if the user wants to use a different command (like, for example, `tofu`), a special flag will need
  to be used to point Deployer to the expected binary.
- Properly configured Terraform Backend, including all authentication if needed. Or Terraform state file
  accessible locally.
- Access to the "new" Deployment Version Terraform code, or Terraform state file accessible locally.
- Access to the "old" Deployment Version Terraform code, or Terraform state file accessible locally.

### Zero-downtime Deployment strategy

1. The first step of performing the Deployment operation will be preparing the "new" Deployment Version
   configuration.

   How it's done is left totally in the hands of the user. User may create a dedicated directory
   and commit it to a Git repository. Or may use some automation mechanism to generate Terraform
   code on the fly from a given template and variables.

   In the documentation, we will provide examples of how this could be done.

1. Once the "new" configuration is ready, we apply it to bring resources up.

   This is being done by calling the `deployer up --target new` command. Deployer commands are described
   in details in a dedicated section below.

   This command is wrapping execution of the `terraform apply` call.

   For commands that wrap execution of `terraform` binary we will make the binary path configurable
   and automatically discovered by default. That will allow users to integrate them with their
   preferred tools (like some other `terraform` wrappers) or easily switch to other implementations
   (like OpenTofu).

1. When resources are created, we connect to the created Runner Managers to verify that the Runner Process
   is up and running.

   This is done by calling the `deployer wait-healthy --target new` command.

   This command will need to read Runner Manager details from Terraform state output. For that, we will
   define a specific name of the output and its structure that users will need to follow.

   Output will be read in two possible ways:
    - by calling `terraform show -json`; in this case we will apply the same rule of making Terraform binary
      configurable
    - by pointing directly to a JSON file with Terraform state.

1. Once the health of the "new" setup is confirmed, we connect to the "old" Runner Managers and initiate
   the Graceful Shutdown.

   This is done by calling the `deployer shutdown --target old` command.

   Optionally, a `deployer shutdown --target old --forceful` command can be used, if the situation
   enforces us to quickly terminate Runner Process without waiting for jobs to be finished first.

   This command reads Terraform state as well.

1. After initializing Runner Process shutdown, we connect to the "old" Runner Managers to verify
   that the Runner Process was terminated.

   This is done by calling the `deployer wait-terminated --target old` command.

   This command reads Terraform state as well.

1. When the "old" Deployment Version is stopped, we can optionally remove the resources.

   This is done by calling the `deployer down --target old` command.

   This command is wrapping execution of the `terraform destroy` call.

   Execution of this step is optional and depends on the strategy chosen by the user.

   Users may want to use a rolling update strategy. In that case, they would like to run this step and
   after it's done - clean up the "old" directory from Git repository or whatever stores the Terraform
   code. Any later change - even if practically it's a revert to the previous configuration - would
   be performed by creating the "new" configuration and executing the Deployment process again.

   Users may want to use a blue/green like update strategy. In that case, they would keep two
   Terraform code directories called "blue" and "green" (or whatever is preferred). When executing
   the Deployment process, one of them would get activated and the second deactivated. Another Deployment
   would switch the states. But the code and resources would be kept, and only the existence of
   Runner Process on the Runner Managers (and, of course, their configuration) would be changed.

### Deployer

Deployer is a binary responsible for executing subsequent steps of the Deployment process.

It wraps around Terraform command calls and reading Terraform state. That could be done, of course,
by directly calling different `terraform ...` commands. But using Deployer for that gives
us some abstraction and limits the Deployment process to use only one command.

Deployer is also the tool that integrates with the Runner Process Wrapper. That allows ensuring
that the Deployment process proceeds only when "new" Deployment Version is truly active and
that the "old" Deployment Version is stopped using the Graceful Shutdown strategy. These both
are requirements for a true Zero-downtime Deployment strategy.

Deployer's code is hosted in GRIT's repository. It makes the Deployer an inseparable part
of GRIT. Yet, as any other module of GRIT, usage of it is totally optional and depends on
users will.

Deployer is a statically linked binary, compiled for every release of GRIT (so also following
GRIT's versioning) and published using GitLab Generic Package Registry, GitLab Releases and
Permanent links to release assets.

Deployer connects to Runner Process Wrapper using the exposed socket. As described
[below](#runner-process-wrapper), Runner Process Wrapper doesn't introduce any authentication
mechanism and defers that to the user. To simplify the design, we expect an [SSH Tunnel](#ssh-connection-strategies)
to be established between Deployer's environment and the Runner Manager. Through that
tunnel a connection to the socket will be made.

To make Deployer's integration easier, any internal logging messages are logged:

- as `NDJSON` - a newline delimited JSON, meaning each log message it's a separate JSON encoded line,
- to `STDERR`, so that Deployer's logging can be easily separated from any other output it could write
  (like, for example, the proxied output of `terraform apply` or `terraform destroy` that's described below).

#### Deployer commands

##### `deployer up`

Syntax: `deployer up --target=<target_name>`

It's a convenient wrapper around the `terraform apply` command executed in a directory accessible through the
`target_name` identifier.

Same outcome could be achieved by directly going to the `target_name` directory and executing `terraform apply`
there.

We provide this command as a simplification, so that user needs to remember and integrate with only the `deployer`
binary. It will also allow use to group commands together to be called by a dedicated alias, which can
be useful in some scenarios.

When executed, the command outputs the `terraform apply` output to `STDOUT`, allowing the user to inspect what
changes were executed. In case of `terraform apply` failure, the command fails as well, returning Terraform's
exit code as its own.

**Exit codes**

| Status  | Value | Description                                                    |
|---------|-------|----------------------------------------------------------------|
| success | 0     | Command succeeded                                              |
| failure | other | Exit code returned by the internal `terraform apply` execution |

##### `deployer wait-healthy`

Syntax: `deployer wait-healthy --target=<target_name> [--retry=<count>] [--timeout=<duration>]`

When called, Deployer reads the Terraform State of the `target_name` Deployment Version. From there
it reads the list of deployed Runner Managers.

It next connects to these Runner Managers and awaits, until dedicated gRPC calls confirm that the
Runner Process was started and is running.

As `deployer wait-healthy` can be called before the resource gets available or the connection tunnel
is fully established, we may need to retry connections to Runner Process Wrapper gRPC server few
times. That part is controlled with the `--retry` flag with a reasonable default count value.

Connection to the Runner Process Wrapper may be established before the Runner Process is fully
started; therefore, gRPC calls can return an "unhealthy" state. The `--timeout` flag with
a reasonable default value may be used to control how long we want to wait before reporting a startup failure.

**Exit codes**

| Status  | Value | Description                                                                                              |
|---------|-------|----------------------------------------------------------------------------------------------------------|
| success | 0     | Command succeeded                                                                                        |
| failure | 1     | Any failure not explicitly covered with a dedicated exit code                                            |
| failure | 2     | Failure due to `--retry` value exceeded when trying to connect the gRPC server                           |
| failure | 3     | Failure due to `--timeout` value exceeded when awaiting for gRPC call to report a running Runner Process |

##### `deployer shutdown`

Syntax: `deployer shutdown --target=<target_name> [--retry=<count>] [--forceful]`

When called, Deployer reads the Terraform State of the `target_name` Deployment Version. From there
it reads the list of deployed Runner Managers.

It next connects to these Runner Managers and initiates the Graceful Shutdown through Runner Process Wrapper
gRPC server.

Similarly as for `wait-healthy`, the `--retry` flag is used to control how many connection attempts to
the gRPC server are allowed before a failure is reported.

`shutdown` can be also called with the `--forceful` flag. In that case, it will initiate a Forceful Shutdown
immediately through the gRPC call.

**Exit codes**

| Status  | Value | Description                                                                    |
|---------|-------|--------------------------------------------------------------------------------|
| success | 0     | Command succeeded                                                              |
| failure | 1     | Any failure not explicitly covered with a dedicated exit code                  |
| failure | 2     | Failure due to `--retry` value exceeded when trying to connect the gRPC server |

##### `deployer wait-terminated`

Syntax: `deployer wait-terminated --target=<target_name> [--retry=<count>] [--timeout=<duration>]`

When called, Deployer reads the Terraform State of the `target_name` Deployment Version. From there
it reads the list of deployed Runner Managers.

It next connects to these Runner Managers and checks for the Runner Process status to be turned into
`terminated`, which means that the shutdown (either Graceful or Forceful) was finished.

Similarly as for `wait-healthy`, the `--retry` flag is used to control how many connection attempts to
the gRPC server are allowed before a failure is reported.

As the Graceful Shutdown may take a long time, the `--timeout` flag can be used to limit how long
Deployer will wait for the rRunner Process to be terminated.

A reasonable default value will be provided. Otherwise, if timeout would be undefined, and Runner Process
would for some reason hang on terminating finished jobs, `deployer wait-terminated` command
would run forever.

**Exit codes**

| Status  | Value | Description                                                                                                         |
|---------|-------|---------------------------------------------------------------------------------------------------------------------|
| success | 0     | Command succeeded                                                                                                   |
| failure | 1     | Any failure not explicitly covered with a dedicated exit code                                                       |
| failure | 2     | Failure due to `--retry` value exceeded when trying to connect the gRPC server                                      |
| failure | 3     | Failure due to `--timeout` value exceeded when awaiting for gRPC call to report the Runner Process to be terminated |

##### `deployer down`

Syntax: `deployer down --target=<target_name>`

It's a convenient wrapper around the `terraform destroy` command executed in a directory accessible through the
`target_name` identifier.

Same outcome could be achieved by directly going to the `target_name` directory and executing `terraform destroy`
there.

We provide this command as a simplification, so that user needs to remember and integrate with only the `deployer`
binary. It will also allow use to group commands together to be called by a dedicated alias, which can
be useful in some scenarios.

When executed, the command outputs the `terraform destroy` output to `STDOUT`, allowing the user to inspect what
changes were executed. In case of `terraform destroy` failure, the command fails as well, returning Terraform's
exit code as its own.

**Exit codes**

| Status  | Value | Description                                                      |
|---------|-------|------------------------------------------------------------------|
| success | 0     | Command succeeded                                                |
| failure | other | Exit code returned by the internal `terraform destroy` execution |

### Runner Process Wrapper

GitLab Runner is a long-running process that can't be just terminated. Termination will affect
all jobs executed at that time with this runner.

Because of that, we have a mechanism in GitLab Runner that is called Graceful Shutdown.
When it's initiated, Runner stops asking for new jobs, but continues execution until
already started jobs are finished. And that may take a lot of time.

Currently, the only way to initiate Graceful Shutdown is to send a `SIGQUIT` to the runner process.
And that's not the best interface to integrate with.

There is also a case of Runner deployed in Kubernetes. There is an easy way to send the signal
to the process on a pod termination. But once the pod is marked for deletion, any Kubernetes
Services connected with that pod are immediately removed. And that affects a possibility
to easily track Runner's metrics exporter during termination.

Process wrapper addresses these concerns. It provides a nicer API (definitely gRPC, we may consider
REST as well for user's manual integrations) and it allows to put the Runner Process into
a graceful shutdown, while still having the main process running.

Wrapper's API should be very basic:

- An endpoint to initiate graceful shutdown.
- An endpoint to enforce forceful shutdown if a situation requires it.
- An endpoint to check the process status.

Wrapper will be implemented in the Runner code, basically becoming a new `gitlab-runner`
binary command. Usage of it will be optional, and it will depend on the user's choice.

When started, it should expose the gRPC server.

Internally it should start the regular `gitlab-runner run` command **as a separate process** that
is fully in control of the Runner Process Wrapper process.

We will not introduce any authentication mechanism for this internal server. We will relay this
responsibility to an external layer. In cases when the server can't be directly exposed
from the Runner node, users should defer to an SSH tunnel or some kind of reverse proxy
with authentication configured for that endpoint.

### SSH Connection strategies

As noted in the [Deployer](#deployer) section, we will use SSH protocol as an interface for
connecting Runner Manager instances and Runner Process Wrapper socket.

Because deployment scenarios may be very different, we will support several strategies
of how this SSH connection should be established.

Once the SSH connection to the Runner Manager is established, a tunneled connection to the
Runner Process Wrapper socket is created, which provides access to Process Wrapper's gRPC
server.

#### SSH Authentication to Runner Manager

For SSH authentication on the Runner Manager hosts, we will enforce the usage of authentication keys.

It's up to the user whether key should be hardcoded on the Runner Manager instance through Terraform,
or whether existing cloud features should be used to inject it into a running instance at the connection time.

GRIT modules [will expose information](#required-terraform-state-output) that will allow Deployer to
establish the SSH connection. In the Authentication scope, this will contain the username and the PEM
encoded value of the authentication private key.

For cases when user will not want to hardcode the SSH key, Deployer will allow to specify the username
and path to the key from command line. If specified, it will take precedence over the values discovered
from the Terraform Output.

#### SSH Strategy: direct

This is the simplest strategy. In this case, Deployer and Runner Manager exist in the same local
network so Deployer can directly open connection to the SSH service on the Runner Manager. It
doesn't require any specific configuration, and directly uses the address discovered from
Terraform Output.

#### SSH Strategy: ProxyJump

A common approach for securing SSH access to the infrastructure is to keep important instances
inside a private network and expose to the public only a server dedicated to entering this network
through SSH. This is commonly referenced as a Bastion server.

In this case, direct access to a `service-a` host is not available. Instead, user needs to SSH
to the `bastion` host and through that host - to `service-a` one, which is available on the
internal network.

While that's easy to do when the user wants to open an interactive shell session, for tunneling
connections that need to be automated. OpenSSH Client, for cases like that, introduces a concept
of `ProxyJump`. With `ProxyJump` user defines the proxy node (so the `bastion` server in our example)
to which OpenSSH Client connects and instructs it to tunnel an SSH connection to the target node.

Deployer will support a similar approach. User will be able to point a `ProxyJump` host through
a dedicated commandline flag, as well as providing username and path to the private key file
that should be used for authenticating with this `ProxyJump` host. This will be dedicated
username/key pair, separated from the username/key pair used for configuring access to the
Runner Manager nodes.

Deployer will establish an SSH connection to the `ProxyJump` node, through it will establish
a tunneled SSH connection to Runner Manager - similarly how it would do it in the `direct`
strategy, and through that - the connection to Runner Process Wrapper socket.

#### SSH Strategy: ProxyCommand

`ProxyJump` approach is a very basic approach. It requires a dedicated host to be exposed
in the "public" network, and it works with SSH service directly.

In some cases, users may want to use a dedicated cloud service that will handle the creation
and authentication of the connection from the "public" network to the "private" network.

An example of this approach is [AWS SSM](https://docs.aws.amazon.com/systems-manager/latest/userguide/ssm-agent.html).

With AWS SSM, users are not managing their own `bastion` servers. Instead, they relay
the responsibility of securing access to internal networks to AWS. After preparing a
specific configuration, including IAM policies changes, in their account, users can
then use a dedicated command of the AWS command line tool.

This command connects to AWS API and authenticates using the standard IAM approach
(so either local credentials, or assumed Service Account, when the request starts within AWS
itself). When authentication succeeds and IAM subsystem confirms that the security policy
grants the requesting entity access to the requested instance, it creates a connection
to this instance's SSH service.

Through that connection, the SSH protocol can perform authentication on the node as it would
be connected directly, while everything is routed through a dedicated AWS networking
and API services.

To support cases like that, OpenSSH Client introduces a concept of `ProxyCommand`. In this case, the
user defines a command that the client will execute (with templating support, to inject some values
dynamically if needed). This command then transfers data through STDIN/STDOUT through which OpenSSH
Client performs SSH protocol communication.

Deployer will support a similar approach. Through a dedicated command line flag, the user will be able
to specify the command that Deployer should execute to establish the connection to the
SSH service on the Runner Manager instance. Through this connection, the SSH library in Deployer will
handle authentication and tunneling to the Runner Process Wrapper socket.

This flag will support templating, so that the `ProxyCommand` can be called with the Runner Manager
node details if needed.

#### SSH Strategy: SSH command

There are cases, where `ProxyCommand` will not work as the infrastructure doesn't expose
SSH service in any way.

An example can be the Google Cloud SSH over IAP feature. In its concept, it's similar to AWS SSM,
but in the details it works slightly differently.

When the user wants to use this to connect to the instance, a dedicated command from Google Cloud command
line tool is used. It also uses the standard `gcloud` authentication mechanisms to authenticate
this call. If it works, `gcloud` connects to the Runner Manager SSH server (through a dedicated
network, owned by Google Cloud, so it requires a specific Firewall rules) and performs all operations
with it. That includes authentication (so it's the `gcloud` command that needs to use the private key)
and opening a tunneled connection to the Runner Process Wrapper socket.

In this case, there is no place for "outer" SSH interaction, like it was in `ProxyCommand` case, where
the command was opening just the connection to the SSH service, but everything else was
left to the OpenSSH Client. Here we get direct access to the Runner Process Wrapper socket.

Deployer will support that approach as well. Through a dedicated command line flag, the user will be
able to specify the command expected to fully establish and manage connection to the Runner Wrapper Process
socket. It will also support templating, similarly to the `ProxyCommand` approach.

### Required Terraform State Output

As noted multiple times above, Deployer will expect that details about Runner Manager and Runner Process
Wrapper will be available through Terraform State Output. For commands that interact with Runner Manager,
Deployer will read the Terraform State (either by calling `terraform show -json` on a given path or by
reading the Terraform State JSON file directly), and in the output section it will seek for the output entry
named `grit_runner_managers`.

Because the configuration using GRIT modules may generate multiple runner managers for the same Deployment
Version (for example, for horizontal scaling or HA purpose), this entry is expected to be a map in the
following structure:

```hcl
type = map(object{
  instance_name   = string
  address         = string
  wrapper_address = string
  username = optional(string, "")
  ssh_key_pem = optional(string, "")
})
```

The key of this map will be used only as an identifier for Deployer's logging. It can use any
arbitrary name that will make it easier to the user to distinguish which Runner Manager a specific
Deplyer operation was related to.

Each Runner Manager entry contains the following five keys:

| Key               | Required | Description                                                                                                                                                                                              |
|-------------------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `instance_name`   | yes      | Name of the instance on which the Runner Manager is executed. This may be useful for users depending on the [ProxyJump](#ssh-strategy-proxyjump) or [SSH Command](#ssh-strategy-ssh-command) strategies. |
| `address`         | yes      | The address of the Runenr Manager host. Either a domain name (in this case Deployer must be able to resolve it) or an IP address through which Runner Manager's SSH service will be accessed.            |
| `wrapper_address` | yes      | The address through which Runner Wrapper Process socket with its gRPC server is accessible locally on the Runner Manager node.                                                                           |
| `username`        | no       | Username through which Runner Manager instance can be accessed, paired with the configured SSH key. If not provided through Terraform Output, user must provide it with command line flag.               |
| `ssh_key_pem`     | no       | PEM encoded SSH private key that should be used to authenticate with the paired username. If not provided through Terraform Output, user must provide it with command line flag.                         |

GRIT modules for AWS and Google should provide this output for the runner managers generated through them. Users
consuming GRIT modules and wanting to use Deployer will be responsible for passing these information from the
modules to the main Terraform Output of their Deployment Version state.

Example output definition:

```hcl
output "grit_runner_managers" {
  value = {
    "runner_manager_1" = {
      "address"         = "runner-manager-1.localdomain:22"
      "username"        = "gitlab-runner"
      "wrapper_address" = "unix:///var/run/gitlab-runner-wrapper.sock"
    }

    "runner_manager_2" = {
      "address"         = "1.2.3.4:22"
      "username"        = "gitlab-runner"
      "wrapper_address" = "unix:///var/run/gitlab-runner-wrapper.sock"
    }

    // ...

    "runner_manager_N" = {
      "address"         = "runner-manager-N.some.local.domain:22"
      "username"        = "example-username"
      "wrapper_address" = "unix:///var/run/wrapper.sock"
    }
  }
}
```

### CI/CD pipeline support

A convenient way to manage Deployment Versions and execute Deployments on them would be to
use the GitOps approach.

> Note that as everything in GRIT, this is an option, not a requirement. Users may decide that
> Deployer binary commands is everything they need, and they want to integrate it into their
> own way of managing Shards configuration. In that case, support for GitOps provided by GRIT
> can be totally ignored.

If GitOps is in use, Terraform code is stored in a Git repository in some structured way, having
dedicated directories for shards and inside them - dedicated directories for Deployment Versions.

A Deployment is initiated by commiting change to the project. That change is next detected
by the CI/CD system, and the Deployment process is executed within the CI/CD pipeline.

Detection of the change depends on the project layout, which as was said, is totally arbitrary.
Users may also have their own design of the CI/CD Pipeline.

Therefore, GRIT's support will be limited to providing CI/CD Components that users can include
in their CI/CD Pipelines. The Components will be responsible for executing the Deployer commands.

While we're sure providing these Components will be a valuable part of GRIT's Zero-downtime Deployment
support, we're not going to focus on that now. We will focus on that part in the second iteration
of working on the Zero-downtime Deployment support.

### Other updates to GRIT

For proper Runner infrastructure management, we need to be able to monitor runner managers.
This design proposes a very dynamic lifetime of the runner manager deployments. For that,
we need to be able to automatically discover them.

We will add a GRIT module for deploying a Prometheus server. It should use Prometheus' automatic
discovery mechanisms to find Runner deployed in the supported clouds.

To allow integrating that monitoring with user's organization monitoring stack, we will make
it possible to - optionally - integrate this Prometheus server with a given
[Mimir](https://grafana.com/oss/mimir/) instance.

## Status

RFC

## Authors

<!-- vale gitlab.Spelling = NO -->

Proposal:

| Role    | Who                                          |
|---------|----------------------------------------------|
| Authors | Tomasz Maczukin, Joseph Burnett, Andy Knight |

Domain experts:

| Area                       | Who             |
|----------------------------|-----------------|
| GitLab Runner              | Tomasz Maczukin |
| Hosted Runners deployments | Tomasz Maczukin |
| GitLab Dedicated           | Andy Knight     |
