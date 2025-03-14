package wait

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
	DefaultTimeout = 30 * time.Minute
)

var (
	errWrapperStatusCheck = errors.New("wrapper status check")
)

type Flags struct {
	Timeout time.Duration
}

//go:generate mockery --name=tfClientMux --inpackage --with-expecter
type tfClientMux interface {
	Execute(ctx context.Context, fn wrapper.CallbackFn) error
}

type loopStatusCheckFn func(ctx context.Context, c wrapper.LoopStatusCheckClient, timeout time.Duration, checkForRunning bool) error

type Service struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags      terraform.Flags
	sshFlags     ssh.Flags
	wrapperFlags wrapper.Flags
	waitFlags    Flags

	tfClientMuxFactory func(*slog.Logger, *terraform.Client, terraform.Flags, ssh.Flags, wrapper.Flags) tfClientMux
	loopStatusCheck    loopStatusCheckFn
}

func New(logger *slog.Logger, tf *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags wrapper.Flags, waitFlags Flags) *Service {
	return &Service{
		logger:       logger.With("tf-flags", tfFlags, "ssh-flags", sshFlags, "wrapper-flags", wrapperFlags, "wait-flags", waitFlags),
		tf:           tf,
		tfFlags:      tfFlags,
		sshFlags:     sshFlags,
		wrapperFlags: wrapperFlags,
		waitFlags:    waitFlags,
		tfClientMuxFactory: func(l *slog.Logger, tf *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags wrapper.Flags) tfClientMux {
			return wrapper.NewMux(l, tf, tfFlags, sshFlags, wrapperFlags)
		},
		loopStatusCheck: wrapper.LoopStatusCheck,
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
	return s.tfClientMuxFactory(s.logger, s.tf, s.tfFlags, s.sshFlags, s.wrapperFlags).
		Execute(ctx, func(ctx context.Context, c wrapper.CallbackClient) error {
			err := s.loopStatusCheck(ctx, c, s.waitFlags.Timeout, checkForRunning)
			if err != nil {
				return fmt.Errorf("%w: %w", errWrapperStatusCheck, err)
			}

			return nil
		})
}
