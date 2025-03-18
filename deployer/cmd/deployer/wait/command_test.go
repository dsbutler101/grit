package wait

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/cli"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/logger"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/services/wait"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/wrapper"
)

func TestCommand_Execute(t *testing.T) {
	testWrapperGRPCConnectionWaitTimeoutExceededError := &wrapper.GRPCConnectionWaitTimeoutExceededError{}
	testWrapperStatusCheckLoopTimeoutExceededError := &wrapper.StatusCheckLoopTimeoutExceededError{}

	waitHealthyPrepareServiceMock := func(err error) func(t *testing.T, sm *mockService, expectedCtx context.Context) {
		return func(t *testing.T, sm *mockService, expectedCtx context.Context) {
			sm.EXPECT().ExecuteWaitHealthy(expectedCtx).Return(err)
		}
	}

	waitTerminatedPrepareServiceMock := func(err error) func(t *testing.T, sm *mockService, expectedCtx context.Context) {
		return func(t *testing.T, sm *mockService, expectedCtx context.Context) {
			sm.EXPECT().ExecuteWaitTerminated(expectedCtx).Return(err)
		}
	}

	assertGRPCConnectionRetryError := func(t *testing.T, err error) {
		assert.ErrorIs(t, err, testWrapperGRPCConnectionWaitTimeoutExceededError)

		var eerr *cli.Error
		if assert.ErrorAs(t, err, &eerr) {
			assert.Equal(t, gRPCConnectionRetryExitCode, eerr.ExitCode())
		}
	}

	assertWrapperStatusCheckLoopTimeoutExceededError := func(t *testing.T, err error) {
		assert.ErrorIs(t, err, testWrapperStatusCheckLoopTimeoutExceededError)

		var eerr *cli.Error
		if assert.ErrorAs(t, err, &eerr) {
			assert.Equal(t, gRPCRunnerProcessReadinessExitCode, eerr.ExitCode())
		}
	}

	assertUnknownFailureError := func(t *testing.T, err error) {
		assert.ErrorIs(t, err, assert.AnError)

		var eerr *cli.Error
		if assert.ErrorAs(t, err, &eerr) {
			assert.Equal(t, unknownFailureExitCode, eerr.ExitCode())
		}
	}

	tests := map[string]struct {
		executionType      executionType
		prepareServiceMock func(t *testing.T, sm *mockService, expectedCtx context.Context)
		assertError        func(t *testing.T, err error)
	}{
		"unknown execution type": {
			executionType:      255,
			prepareServiceMock: func(_ *testing.T, _ *mockService, _ context.Context) {},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errUnknownExecutionType)
			},
		},
		"service execution for wait healthy success": {
			executionType:      executionTypeHealthyCheck,
			prepareServiceMock: waitHealthyPrepareServiceMock(nil),
		},
		"service execution for wait terminated success": {
			executionType:      executionTypeTerminatedCheck,
			prepareServiceMock: waitTerminatedPrepareServiceMock(nil),
		},
		"service execution for wait healthy fails with GRPCConnectionWaitTimeoutExceededError": {
			executionType:      executionTypeHealthyCheck,
			prepareServiceMock: waitHealthyPrepareServiceMock(testWrapperGRPCConnectionWaitTimeoutExceededError),
			assertError:        assertGRPCConnectionRetryError,
		},
		"service execution for wait terminated fails with GRPCConnectionWaitTimeoutExceededError": {
			executionType:      executionTypeTerminatedCheck,
			prepareServiceMock: waitTerminatedPrepareServiceMock(testWrapperGRPCConnectionWaitTimeoutExceededError),
			assertError:        assertGRPCConnectionRetryError,
		},
		"service execution for wait healthy fails with StatusCheckLoopTimeoutExceededError": {
			executionType:      executionTypeHealthyCheck,
			prepareServiceMock: waitHealthyPrepareServiceMock(testWrapperStatusCheckLoopTimeoutExceededError),
			assertError:        assertWrapperStatusCheckLoopTimeoutExceededError,
		},
		"service execution for wait terminated fails with StatusCheckLoopTimeoutExceededError": {
			executionType:      executionTypeTerminatedCheck,
			prepareServiceMock: waitTerminatedPrepareServiceMock(testWrapperStatusCheckLoopTimeoutExceededError),
			assertError:        assertWrapperStatusCheckLoopTimeoutExceededError,
		},
		"service execution for wait healthy fails with unknown failure": {
			executionType:      executionTypeHealthyCheck,
			prepareServiceMock: waitHealthyPrepareServiceMock(assert.AnError),
			assertError:        assertUnknownFailureError,
		},
		"service execution for wait terminated fails with unknown failure": {
			executionType:      executionTypeTerminatedCheck,
			prepareServiceMock: waitTerminatedPrepareServiceMock(assert.AnError),
			assertError:        assertUnknownFailureError,
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			require.NotNil(t, tt.prepareServiceMock, "prepareServiceMock must be defined in test definition")

			sMock := newMockService(t)
			tt.prepareServiceMock(t, sMock, ctx)

			testTFClient := new(terraform.Client)

			c := newCmd(logger.New(), testTFClient, tt.executionType, func(s *slog.Logger, tfClient *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags wrapper.Flags, waitFlags wait.Flags) service {
				assert.Equal(t, testTFClient, tfClient)
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
