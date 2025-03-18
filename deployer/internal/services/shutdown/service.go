package shutdown

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/wrapper"
)

const (
	DefaultRetry = 3
)

var (
	errInitiatingGracefulShutdown = errors.New("initiating graceful shutdown")
	errInitiatingForcefulShutdown = errors.New("initiating forceful shutdown")
)

type Flags struct {
	Forceful bool
}

//go:generate mockery --name=tfClientMux --inpackage --with-expecter
type tfClientMux interface {
	Execute(ctx context.Context, fn wrapper.CallbackFn) error
}

type loopStatusCheckFn func(ctx context.Context, c wrapper.LoopStatusCheckClient, timeout time.Duration, checkForRunning bool) error

type Service struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags       terraform.Flags
	sshFlags      ssh.Flags
	wrapperFlags  wrapper.Flags
	shutdownFlags Flags

	tfClientMuxFactory func(*slog.Logger, *terraform.Client, terraform.Flags, ssh.Flags, wrapper.Flags) tfClientMux
}

func New(logger *slog.Logger, tf *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags wrapper.Flags, shutdownFlags Flags) *Service {
	return &Service{
		logger:        logger.With("operation", "shutdown", "tf-flags", tfFlags, "ssh-flags", sshFlags, "wrapper-flags", wrapperFlags, "shutdown-flags", shutdownFlags),
		tf:            tf,
		tfFlags:       tfFlags,
		sshFlags:      sshFlags,
		wrapperFlags:  wrapperFlags,
		shutdownFlags: shutdownFlags,
		tfClientMuxFactory: func(l *slog.Logger, tf *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags wrapper.Flags) tfClientMux {
			return wrapper.NewMux(l, tf, tfFlags, sshFlags, wrapperFlags)
		},
	}
}

func (s *Service) Execute(ctx context.Context) error {
	return s.tfClientMuxFactory(s.logger, s.tf, s.tfFlags, s.sshFlags, s.wrapperFlags).
		Execute(ctx, func(ctx context.Context, c wrapper.CallbackClient) error {
			if s.shutdownFlags.Forceful {
				err := c.InitForcefulShutdown(ctx)
				if err != nil {
					return fmt.Errorf("%w: %w", errInitiatingForcefulShutdown, err)
				}
				return nil
			}

			err := c.InitGracefulShutdown(ctx)
			if err != nil {
				return fmt.Errorf("%w: %w", errInitiatingGracefulShutdown, err)
			}

			return nil
		})
}
