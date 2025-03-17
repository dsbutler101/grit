package shutdown

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/cli"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/logger"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/services/shutdown"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/wrapper"
)

func TestCommand_Execute(t *testing.T) {
	tests := map[string]struct {
		serviceError error
		assertError  func(t *testing.T, err error)
	}{
		"service execution success": {},
		"service execution fails with unknown error": {
			serviceError: assert.AnError,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)

				var eerr *cli.Error
				if assert.ErrorAs(t, err, &eerr) {
					assert.Equal(t, unknownFailureExitCode, eerr.ExitCode())
				}
			},
		},
		"service execution fails with GRPCConnectionWaitTimeoutExceededError": {
			serviceError: &wrapper.GRPCConnectionWaitTimeoutExceededError{},
			assertError: func(t *testing.T, err error) {
				var eerr *cli.Error
				if assert.ErrorAs(t, err, &eerr) {
					assert.Equal(t, gRPCConnectionRetryExitCode, eerr.ExitCode())
				}
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			testTFClient := new(terraform.Client)

			sMock := newMockService(t)
			sMock.EXPECT().Execute(ctx).Return(tt.serviceError)

			c := newCmd(logger.New(), testTFClient, func(logger *slog.Logger, client *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags wrapper.Flags, shutdownFlags shutdown.Flags) service {
				assert.Equal(t, testTFClient, client)
				return sMock
			})

			err := c.Execute(ctx, nil, nil)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
