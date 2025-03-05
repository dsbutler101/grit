package down

import (
	"context"
	"errors"
	"log/slog"

	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/cmd/deployer/base"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/cli"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/services/down"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
)

const (
	unknownFailureExitCode = 1
)

type cmd struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags *terraform.Flags
}

func (c *cmd) Execute(ctx context.Context, _ *cobra.Command, _ []string) error {
	err := down.New(c.logger, c.tf, *c.tfFlags).Execute(ctx)
	if err == nil {
		return nil
	}

	var terr *terraform.CommandError
	if errors.As(err, &terr) {
		return cli.NewError(terr.ExitCode(), err)
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
		Use:     "down",
		Short:   "Removes Deployment Version through Terraform",
		PreRunE: func(_ *cobra.Command, _ []string) error {
			err := c.tfFlags.Validate()
			if err != nil {
				return err
			}

			return nil
		},
		RunE: cli.BuildRunEFromCommandExecutor(c),
	}

	c.tfFlags = base.SetupTFFlagsTargetOnly(cc)

	return cc
}
