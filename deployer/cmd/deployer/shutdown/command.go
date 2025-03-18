package shutdown

import (
	"context"
	"errors"
	"log/slog"

	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/cmd/deployer/base"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/cli"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/services/shutdown"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/wrapper"
)

const (
	unknownFailureExitCode      = 1
	gRPCConnectionRetryExitCode = 2
)

//go:generate mockery --name=service --inpackage --with-expecter
type service interface {
	Execute(context.Context) error
}

type serviceFactoryFn func(*slog.Logger, *terraform.Client, terraform.Flags, ssh.Flags, wrapper.Flags, shutdown.Flags) service

type cmd struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags       *terraform.Flags
	sshFlags      *ssh.Flags
	wrapperFlags  *wrapper.Flags
	shutdownFlags *shutdown.Flags

	serviceFactory serviceFactoryFn
}

func newCmd(logger *slog.Logger, tf *terraform.Client, sf serviceFactoryFn) *cmd {
	return &cmd{
		logger:         logger,
		tf:             tf,
		tfFlags:        new(terraform.Flags),
		wrapperFlags:   new(wrapper.Flags),
		sshFlags:       new(ssh.Flags),
		shutdownFlags:  new(shutdown.Flags),
		serviceFactory: sf,
	}
}

func (c *cmd) Execute(ctx context.Context, _ *cobra.Command, _ []string) error {
	err := c.serviceFactory(c.logger, c.tf, *c.tfFlags, *c.sshFlags, *c.wrapperFlags, *c.shutdownFlags).Execute(ctx)
	if err == nil {
		return nil
	}

	var rerr *wrapper.GRPCConnectionWaitTimeoutExceededError
	if errors.As(err, &rerr) {
		return cli.NewError(gRPCConnectionRetryExitCode, err)
	}

	return cli.NewError(unknownFailureExitCode, err)
}

func New(logger *slog.Logger, tf *terraform.Client, cmdGroup cobra.Group) *cobra.Command {
	c := newCmd(logger, tf, func(logger *slog.Logger, client *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags wrapper.Flags, shutdownFlags shutdown.Flags) service {
		return shutdown.New(logger, client, tfFlags, sshFlags, wrapperFlags, shutdownFlags)
	})

	cc := &cobra.Command{
		GroupID: cmdGroup.ID,
		Use:     "shutdown",
		Short:   "Initiates Runner Managers shutdown on the given Deployment Version of the Shard",
		Long:    "Connects to Runner Manager's Runner Process Wrapper's gRPC server and initiates the Runner Process shutdown",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			err := c.tfFlags.Validate()
			if err != nil {
				return err
			}

			err = c.sshFlags.Validate()
			if err != nil {
				return err
			}

			return nil
		},
		RunE: cli.BuildRunEFromCommandExecutor(c),
	}

	base.SetupAllTFFlags(cc, c.tfFlags)
	base.SetupWrapperFlags(cc, c.wrapperFlags)
	base.SetupSSHFlags(cc, c.sshFlags)

	cc.PersistentFlags().BoolVar(&c.shutdownFlags.Forceful, "forceful", false, "Initiate Forceful Shutdown instead of Graceful Shutdown")

	return cc
}
