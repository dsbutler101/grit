package wait

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/wrapper"
)

const (
	DefaultTimeout = 30 * time.Minute
)

type Flags struct {
	Timeout time.Duration
}

type Service struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags      terraform.Flags
	sshFlags     ssh.Flags
	wrapperFlags wrapper.Flags
	waitFlags    Flags
}

func New(logger *slog.Logger, tf *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags wrapper.Flags, waitFlags Flags) *Service {
	return &Service{
		logger:       logger.With("tf-flags", tfFlags, "ssh-flags", sshFlags, "wrapper-flags", wrapperFlags, "wait-flags", waitFlags),
		tf:           tf,
		tfFlags:      tfFlags,
		sshFlags:     sshFlags,
		wrapperFlags: wrapperFlags,
		waitFlags:    waitFlags,
	}
}

func (s *Service) ExecuteWaitHealthy(ctx context.Context) error {
	s.logger = s.logger.With("operation", "wait-healthy")

	return s.execute(ctx, wrapper.CheckForRunning)
}

func (s *Service) ExecuteWaitTerminated(ctx context.Context) error {
	s.logger = s.logger.With("operation", "wait-terminated")

	return s.execute(ctx, wrapper.CheckForStopped)
}

func (s *Service) execute(ctx context.Context, checkForRunning bool) error {
	return wrapper.NewMux(s.logger, s.tf, s.tfFlags, s.sshFlags, s.wrapperFlags).
		Execute(ctx, func(ctx context.Context, c *wrapper.Client) error {
			err := wrapper.LoopStatusCheck(ctx, c, s.waitFlags.Timeout, checkForRunning)
			if err != nil {
				return fmt.Errorf("wrapper status check: %w", err)
			}

			return nil
		})
}
