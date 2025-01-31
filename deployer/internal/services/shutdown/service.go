package shutdown

import (
	"context"
	"fmt"
	"log/slog"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/wrapper"
)

const (
	DefaultRetry = 3
)

type Flags struct {
	Forceful bool
}

type Service struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags       terraform.Flags
	sshFlags      ssh.Flags
	wrapperFlags  wrapper.Flags
	shutdownFlags Flags
}

func New(logger *slog.Logger, tf *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags wrapper.Flags, shutdownFlags Flags) *Service {
	return &Service{
		logger:        logger.With("operation", "shutdown", "tf-flags", tfFlags, "ssh-flags", sshFlags, "wrapper-flags", wrapperFlags, "shutdown-flags", shutdownFlags),
		tf:            tf,
		tfFlags:       tfFlags,
		sshFlags:      sshFlags,
		wrapperFlags:  wrapperFlags,
		shutdownFlags: shutdownFlags,
	}
}

func (s *Service) Execute(ctx context.Context) error {
	return wrapper.NewMux(s.logger, s.tf, s.tfFlags, s.sshFlags, s.wrapperFlags).
		Execute(ctx, func(ctx context.Context, c *wrapper.Client) error {
			if s.shutdownFlags.Forceful {
				err := c.InitForcefulShutdown(ctx)
				if err != nil {
					return fmt.Errorf("initiating forceful shutdown: %w", err)
				}
				return nil
			}

			err := c.InitGracefulShutdown(ctx)
			if err != nil {
				return fmt.Errorf("initiating graceful shutdown: %w", err)
			}

			return nil
		})
}
