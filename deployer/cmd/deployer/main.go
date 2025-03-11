package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/cmd/deployer/down"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/cmd/deployer/shutdown"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/cmd/deployer/up"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/cmd/deployer/version"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/cmd/deployer/wait"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/cli"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/logger"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
)

const (
	successfulExitCode   = 0
	unknownErrorExitCode = 1
)

var (
	logL *slog.LevelVar
	log  *slog.Logger

	tf *terraform.Client

	logDebug   bool
	tfExecPath string
)

func main() {
	ctx, cancelFn := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelFn()

	setupLogger()
	setupTFClient()
	rootCmd := setupRootCMD()

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(determineExitCode(err))
	}

	os.Exit(successfulExitCode)
}

func setupLogger() {
	logL = &slog.LevelVar{}
	logL.Set(slog.LevelInfo)

	logOpts := []logger.Option{
		logger.WithLeveler(logL),
		logger.WithAddSource(deployer.AddSources()),
		logger.WithCustomLogFormat(deployer.CustomLogFormat()),
	}

	log = logger.New(logOpts...)
}

func setupTFClient() {
	var err error

	tf, err = terraform.New(log)
	if err != nil {
		log.Error("Could not create Terraform client", "error", err)
		os.Exit(unknownErrorExitCode)
	}
}

func setupRootCMD() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "deployer",
		Short:         "GRIT Zero-downtime deployment tool",
		Version:       deployer.VersionInfo().ExtendedString(),
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if logDebug {
				logL.Set(slog.LevelDebug)
			}

			if tf.ExecPath() != tfExecPath {
				tf.SetExecPath(tfExecPath)
			}

			log.With("pid", os.Getpid()).Info("Starting deployer")
		},
	}

	rootCmd.PersistentFlags().BoolVar(&logDebug, "debug", false, "Set log level to debug")
	rootCmd.PersistentFlags().StringVar(&tfExecPath, "tf-exec-path", tf.ExecPath(), "Path to Terraform executable")
	tfGroup := cobra.Group{ID: "tf", Title: "Terraform maintenance"}

	wrapperGroup := cobra.Group{ID: "wrapper", Title: "Runner process wrapper integration"}
	rootCmd.AddGroup(&tfGroup, &wrapperGroup)

	for _, cmd := range []*cobra.Command{
		up.New(log, tf, tfGroup),
		down.New(log, tf, tfGroup),
		shutdown.New(log, tf, wrapperGroup),
		wait.NewHealthy(log, tf, wrapperGroup),
		wait.NewTerminated(log, tf, wrapperGroup),
		version.New(),
	} {
		rootCmd.AddCommand(cmd)
	}

	return rootCmd
}

func determineExitCode(err error) int {
	exitCode := unknownErrorExitCode

	if errors.Is(err, context.Canceled) {
		log.Error("termination signal received")
		return exitCode
	}

	var cliErr *cli.Error
	if errors.As(err, &cliErr) {
		exitCode = cliErr.ExitCode()
	}

	log.Error("failed to execute command", "error", err, "exitCode", exitCode)

	return exitCode
}
