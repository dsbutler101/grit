package tfexec

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/cli"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/logger"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
)

func TestCommand_Execute(t *testing.T) {
	testTFExitCode := 123

	upPrepareServiceMock := func(err error) func(t *testing.T, sm *mockService, expectedCtx context.Context) {
		return func(t *testing.T, sm *mockService, expectedCtx context.Context) {
			sm.EXPECT().ExecuteUp(expectedCtx).Return(err)
		}
	}

	downPrepareServiceMock := func(err error) func(t *testing.T, sm *mockService, expectedCtx context.Context) {
		return func(t *testing.T, sm *mockService, expectedCtx context.Context) {
			sm.EXPECT().ExecuteDown(expectedCtx).Return(err)
		}
	}

	assertTerraformCommandError := func(t *testing.T, err error) {
		assert.ErrorIs(t, err, assert.AnError)

		var eerr *cli.Error
		if assert.ErrorAs(t, err, &eerr) {
			assert.Equal(t, testTFExitCode, eerr.ExitCode())
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
		"service execution for up success": {
			executionType:      executionTypeUp,
			prepareServiceMock: upPrepareServiceMock(nil),
		},
		"service execution for down success": {
			executionType:      executionTypeDown,
			prepareServiceMock: downPrepareServiceMock(nil),
		},
		"service execution for up fails with a terraform command error": {
			executionType:      executionTypeUp,
			prepareServiceMock: upPrepareServiceMock(terraform.NewCommandError("test-error", testTFExitCode, assert.AnError)),
			assertError:        assertTerraformCommandError,
		},
		"service execution for down fails with a terraform command error": {
			executionType:      executionTypeDown,
			prepareServiceMock: downPrepareServiceMock(terraform.NewCommandError("test-error", testTFExitCode, assert.AnError)),
			assertError:        assertTerraformCommandError,
		},

		"service execution for up fails with unknown error": {
			executionType:      executionTypeUp,
			prepareServiceMock: upPrepareServiceMock(assert.AnError),
			assertError:        assertUnknownFailureError,
		},
		"service execution for down fails with unknown error": {
			executionType:      executionTypeDown,
			prepareServiceMock: downPrepareServiceMock(assert.AnError),
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

			c := newCmd(logger.New(), testTFClient, tt.executionType, func(s *slog.Logger, tfClient *terraform.Client, flags terraform.Flags) service {
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
