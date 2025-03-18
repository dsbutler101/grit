package tfexec

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
)

var (
	errInitializingTerraform        = errors.New("initializing terraform ")
	errApplyingTerraformResources   = errors.New("applying terraform resources")
	errDestroyingTerraformResources = errors.New("destroying terraform resources")
)

//go:generate mockery --name=tfClient --inpackage --with-expecter
type tfClient interface {
	Init(ctx context.Context, dir string) error
	Apply(ctx context.Context, dir string) error
	Destroy(ctx context.Context, dir string) error
}

type Service struct {
	logger *slog.Logger
	tf     tfClient

	tfFlags terraform.Flags
}

func New(logger *slog.Logger, tf tfClient, tfFlags terraform.Flags) *Service {
	return &Service{
		logger:  logger.With("operation", "down", "tf-flags", tfFlags),
		tf:      tf,
		tfFlags: tfFlags,
	}
}

func (s *Service) ExecuteUp(ctx context.Context) error {
	err := s.tf.Init(ctx, s.tfFlags.Target)
	if err != nil {
		return fmt.Errorf("%w: %w", errInitializingTerraform, err)
	}

	err = s.tf.Apply(ctx, s.tfFlags.Target)
	if err != nil {
		return fmt.Errorf("%w: %w", errApplyingTerraformResources, err)
	}

	return nil
}

func (s *Service) ExecuteDown(ctx context.Context) error {
	err := s.tf.Init(ctx, s.tfFlags.Target)
	if err != nil {
		return fmt.Errorf("%w: %w", errInitializingTerraform, err)
	}

	err = s.tf.Destroy(ctx, s.tfFlags.Target)
	if err != nil {
		return fmt.Errorf("%w: %w", errDestroyingTerraformResources, err)
	}

	return nil
}
