package tfexec

import (
	"context"
	"errors"
	"log/slog"

	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/cmd/deployer/base"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/cli"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/services/down"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/services/up"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
)

const (
	unknownFailureExitCode = 1
)

//go:generate mockery --name=service --inpackage --with-expecter
type service interface {
	Execute(context.Context) error
}

type serviceFactoryFn func(*slog.Logger, *terraform.Client, terraform.Flags) service

type cmd struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags *terraform.Flags

	serviceFactory serviceFactoryFn
}

func newCmd(logger *slog.Logger, tf *terraform.Client, sf serviceFactoryFn) *cmd {
	return &cmd{
		logger:         logger,
		tf:             tf,
		tfFlags:        new(terraform.Flags),
		serviceFactory: sf,
	}
}

func (c *cmd) Execute(ctx context.Context, _ *cobra.Command, _ []string) error {
	err := c.serviceFactory(c.logger, c.tf, *c.tfFlags).Execute(ctx)
	if err == nil {
		return nil
	}

	var terr *terraform.CommandError
	if errors.As(err, &terr) {
		return cli.NewError(terr.ExitCode(), err)
	}

	return cli.NewError(unknownFailureExitCode, err)
}

func NewDown(logger *slog.Logger, tf *terraform.Client, cmdGroup cobra.Group) *cobra.Command {
	cc := newCobraCmd(logger, tf, cmdGroup, func(logger *slog.Logger, client *terraform.Client, flags terraform.Flags) service {
		return down.New(logger, client, flags)
	})

	cc.Use = "down"
	cc.Short = "Removes Deployment Version through Terraform"

	return cc
}

func NewUp(logger *slog.Logger, tf *terraform.Client, cmdGroup cobra.Group) *cobra.Command {
	cc := newCobraCmd(logger, tf, cmdGroup, func(logger *slog.Logger, client *terraform.Client, flags terraform.Flags) service {
		return up.New(logger, client, flags)
	})

	cc.Use = "up"
	cc.Short = "Brings up Deployment Version through Terraform"

	return cc
}

func newCobraCmd(logger *slog.Logger, tf *terraform.Client, cmdGroup cobra.Group, sf serviceFactoryFn) *cobra.Command {
	c := newCmd(logger, tf, sf)

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

	base.SetupTFFlagsTargetOnly(cc, c.tfFlags)

	return cc
}
