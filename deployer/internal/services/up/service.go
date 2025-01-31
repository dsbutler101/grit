package up

import (
	"context"
	"fmt"
	"log/slog"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
)

type Service struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags terraform.Flags
}

func New(logger *slog.Logger, tf *terraform.Client, tfFlags terraform.Flags) *Service {
	return &Service{
		logger:  logger.With("operation", "up", "tf-flags", tfFlags),
		tf:      tf,
		tfFlags: tfFlags,
	}
}

func (s *Service) Execute(ctx context.Context) error {
	err := s.tf.Init(ctx, s.tfFlags.Target)
	if err != nil {
		return fmt.Errorf("initializing terraform code: %w", err)
	}

	err = s.tf.Apply(ctx, s.tfFlags.Target)
	if err != nil {
		return fmt.Errorf("applying terraform code: %w", err)
	}

	return nil
}
