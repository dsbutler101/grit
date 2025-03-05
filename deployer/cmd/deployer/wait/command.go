package wait

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/cmd/deployer/base"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/cli"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/services/wait"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/wrapper"
)

const (
	unknownFailureExitCode             = 1
	gRPCConnectionRetryExitCode        = 2
	gRPCRunnerProcessReadinessExitCode = 3
)

type executionType uint8

const (
	executionTypeHealthyCheck executionType = iota
	executionTypeTerminatedCheck
)

type cmd struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags      *terraform.Flags
	wrapperFlags *wrapper.Flags
	sshFlags     *ssh.Flags
	waitFlags    *wait.Flags

	executionType executionType
}

func (c *cmd) Execute(ctx context.Context, _ *cobra.Command, _ []string) error {
	svc := wait.New(c.logger, c.tf, *c.tfFlags, *c.sshFlags, *c.wrapperFlags, *c.waitFlags)

	executionsMap := map[executionType]func(ctx context.Context) error{
		executionTypeHealthyCheck:    svc.ExecuteWaitHealthy,
		executionTypeTerminatedCheck: svc.ExecuteWaitTerminated,
	}

	execute, ok := executionsMap[c.executionType]
	if !ok {
		return fmt.Errorf("unknown execution type: %v", c.executionType)
	}

	err := execute(ctx)
	if err == nil {
		return nil
	}

	var rerr *wrapper.GRPCConnectionWaitTimeoutExceededError
	if errors.As(err, &rerr) {
		return cli.NewError(gRPCConnectionRetryExitCode, err)
	}

	var terr *wrapper.StatusCheckLoopTimeoutExceededError
	if errors.As(err, &terr) {
		return cli.NewError(gRPCRunnerProcessReadinessExitCode, terr)
	}

	return cli.NewError(unknownFailureExitCode, err)
}

func NewHealthy(logger *slog.Logger, tf *terraform.Client, cmdGroup cobra.Group) *cobra.Command {
	cc := newCmd(logger, tf, cmdGroup, executionTypeHealthyCheck)

	cc.Use = "wait-healthy"
	cc.Short = "Awaits Runner Managers startup for the given Deployment Version of the Shard"
	cc.Long = "Connects to Runner Manager's Runner Process Wrapper's gRPC server and waits until it confirms that Runner Process is ready"

	return cc
}

func NewTerminated(logger *slog.Logger, tf *terraform.Client, cmdGroup cobra.Group) *cobra.Command {
	cc := newCmd(logger, tf, cmdGroup, executionTypeTerminatedCheck)

	cc.Use = "wait-terminated"
	cc.Short = "Awaits Runner Managers termination for the given Deployment Version of the Shard"
	cc.Long = "Connects to Runner Manager's Runner Process Wrapper's gRPC server and waits until it confirms that Runner Process was terminated"

	return cc
}

func newCmd(logger *slog.Logger, tf *terraform.Client, cmdGroup cobra.Group, et executionType) *cobra.Command {
	c := &cmd{
		logger:        logger,
		tf:            tf,
		executionType: et,
	}

	cc := &cobra.Command{
		GroupID: cmdGroup.ID,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			err := c.tfFlags.Validate()
			if err != nil {
				return err
			}

			return nil
		},
		RunE: cli.BuildRunEFromCommandExecutor(c),
	}

	c.tfFlags = base.SetupAllTFFlags(cc)
	c.wrapperFlags = base.SetupWrapperFlags(cc)
	c.sshFlags = base.SetupSSHFlags(cc)
	c.waitFlags = base.SetupWaitFlags(cc)

	return cc
}
