package cli

import (
	"context"

	"github.com/spf13/cobra"
)

type CommandExecutor interface {
	Execute(ctx context.Context, cmd *cobra.Command, args []string) error
}

func BuildCommandExecutor(c CommandExecutor) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return c.Execute(cmd.Context(), cmd, args)
	}
}
