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

type cmd struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags       *terraform.Flags
	sshFlags      *ssh.Flags
	wrapperFlags  *wrapper.Flags
	shutdownFlags *shutdown.Flags
}

func (c *cmd) Execute(ctx context.Context, _ *cobra.Command, _ []string) error {
	err := shutdown.New(c.logger, c.tf, *c.tfFlags, *c.sshFlags, *c.wrapperFlags, *c.shutdownFlags).Execute(ctx)
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
	c := &cmd{
		logger: logger,
		tf:     tf,
	}

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

			return nil
		},
		RunE: cli.BuildCommandExecutor(c),
	}

	c.tfFlags = base.SetupAllTFFlags(cc)
	c.wrapperFlags = base.SetupWrapperFlags(cc)
	c.sshFlags = base.SetupSSHFlags(cc)
	c.shutdownFlags = &shutdown.Flags{}

	cc.PersistentFlags().BoolVar(&c.shutdownFlags.Forceful, "forceful", false, "Initiate Forceful Shutdown instead of Graceful Shutdown")

	return cc
}
