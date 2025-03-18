package tfexec

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/cmd/deployer/base"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/cli"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/services/tfexec"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
)

const (
	unknownFailureExitCode = 1
)

type executionType uint8

const (
	executionTypeUp executionType = iota
	executionTypeDown
)

var (
	errUnknownExecutionType = errors.New("unknown execution type")
)

//go:generate mockery --name=service --inpackage --with-expecter
type service interface {
	ExecuteUp(context.Context) error
	ExecuteDown(context.Context) error
}

type serviceFactoryFn func(*slog.Logger, *terraform.Client, terraform.Flags) service

type cmd struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags *terraform.Flags

	executionType executionType

	serviceFactory serviceFactoryFn
}

func newCmd(logger *slog.Logger, tf *terraform.Client, et executionType, sf serviceFactoryFn) *cmd {
	return &cmd{
		logger:         logger,
		tf:             tf,
		tfFlags:        new(terraform.Flags),
		executionType:  et,
		serviceFactory: sf,
	}
}

func (c *cmd) Execute(ctx context.Context, _ *cobra.Command, _ []string) error {
	svc := c.serviceFactory(c.logger, c.tf, *c.tfFlags)

	executionsMap := map[executionType]func(ctx context.Context) error{
		executionTypeUp:   svc.ExecuteUp,
		executionTypeDown: svc.ExecuteDown,
	}

	execute, ok := executionsMap[c.executionType]
	if !ok {
		return fmt.Errorf("%w: %v", errUnknownExecutionType, c.executionType)
	}

	err := execute(ctx)
	if err == nil {
		return nil
	}

	var terr *terraform.CommandError
	if errors.As(err, &terr) {
		return cli.NewError(terr.ExitCode(), err)
	}

	return cli.NewError(unknownFailureExitCode, err)
}

func NewUp(logger *slog.Logger, tf *terraform.Client, cmdGroup cobra.Group) *cobra.Command {
	cc := newCobraCmd(logger, tf, cmdGroup, executionTypeUp)

	cc.Use = "up"
	cc.Short = "Brings up Deployment Version through Terraform"

	return cc
}

func NewDown(logger *slog.Logger, tf *terraform.Client, cmdGroup cobra.Group) *cobra.Command {
	cc := newCobraCmd(logger, tf, cmdGroup, executionTypeDown)

	cc.Use = "down"
	cc.Short = "Removes Deployment Version through Terraform"

	return cc
}

func newCobraCmd(logger *slog.Logger, tf *terraform.Client, cmdGroup cobra.Group, et executionType) *cobra.Command {
	c := newCmd(logger, tf, et, func(logger *slog.Logger, client *terraform.Client, flags terraform.Flags) service {
		return tfexec.New(logger, client, flags)
	})

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
