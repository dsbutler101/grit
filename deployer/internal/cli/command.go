package cli

import (
	"context"

	"github.com/spf13/cobra"
)

//go:generate mockery --name=CommandExecutor --inpackage --with-expecter
type CommandExecutor interface {
	Execute(ctx context.Context, cmd *cobra.Command, args []string) error
}

func BuildRunEFromCommandExecutor(c CommandExecutor) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return c.Execute(cmd.Context(), cmd, args)
	}
}
