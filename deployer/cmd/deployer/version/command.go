package version

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/cli"
)

type cmd struct {
	outputJSON bool
}

func (c *cmd) Execute(ctx context.Context, cmd *cobra.Command, args []string) error {
	if c.outputJSON {
		ver, err := deployer.VersionInfo().JSON()
		if err != nil {
			return err
		}

		fmt.Println(ver)
		return nil
	}

	fmt.Println(deployer.VersionInfo().DetailedString())

	return nil
}

func New() *cobra.Command {
	c := &cmd{}

	cc := &cobra.Command{
		Use:   "version",
		Short: "Prints Deployer version and exits",
		RunE:  cli.BuildCommandExecutor(c),
	}

	cc.PersistentFlags().BoolVarP(&c.outputJSON, "json", "j", false, "Print out in JSON format")

	return cc
}
