package wait

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/logger"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/wrapper"
)

func TestService_ExecuteWaitHealthy(t *testing.T) {
	testServiceExecute(t, wrapper.CheckForRunning, func(ctx context.Context, s *Service) error {
		return s.ExecuteWaitHealthy(ctx)
	})
}

func testServiceExecute(t *testing.T, expectedCheckForRunning bool, testFn func(ctx context.Context, s *Service) error) {
	testTimeout := 1 * time.Millisecond
	testTimeout2 := 10 * time.Millisecond

	testTFFlags := terraform.Flags{
		Target: "test-target",
	}
	testSSHFlags := ssh.Flags{
		Username: "test-username",
		KeyFile:  "test-key-file",
	}
	testWrapperFlags := wrapper.Flags{
		ConnectionTimeout: testTimeout,
	}
	testFlags := Flags{
		Timeout: testTimeout2,
	}

	noopMockLoopStatusCheck := func(t *testing.T, expectedCtx context.Context, expectedTimeout time.Duration, expectedCheckForRunning bool) loopStatusCheckFn {
		return func(ctx context.Context, _ wrapper.LoopStatusCheckClient, timeout time.Duration, checkForRunning bool) error {
			assert.Equal(t, expectedCtx, ctx)
			assert.Equal(t, expectedTimeout, timeout)
			assert.Equal(t, expectedCheckForRunning, checkForRunning)

			return nil
		}
	}

	defaultPrepareTFClientMuxMock := func(t *testing.T, m *mockTfClientMux, expectedCtx context.Context) {
		m.EXPECT().Execute(expectedCtx, mock.Anything).RunAndReturn(func(ctx context.Context, fn wrapper.CallbackFn) error {
			return fn(ctx, nil)
		})
	}

	tests := map[string]struct {
		tfFlags                terraform.Flags
		sshFlags               ssh.Flags
		wrapperFlags           wrapper.Flags
		flags                  Flags
		prepareTFClientMuxMock func(t *testing.T, m *mockTfClientMux, expectedCtx context.Context)
		mockLoopStatusCheck    func(t *testing.T, expectedCtx context.Context, expectedTimeout time.Duration, expectedCheckForRunning bool) loopStatusCheckFn
		assertError            func(t *testing.T, err error)
	}{
		"tfClientMux execution error": {
			tfFlags:      testTFFlags,
			sshFlags:     testSSHFlags,
			wrapperFlags: testWrapperFlags,
			flags:        testFlags,
			prepareTFClientMuxMock: func(t *testing.T, m *mockTfClientMux, expectedCtx context.Context) {
				m.EXPECT().Execute(expectedCtx, mock.Anything).Return(assert.AnError)
			},
			mockLoopStatusCheck: noopMockLoopStatusCheck,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		"tfClientMux loopStatusCheck error": {
			tfFlags:                testTFFlags,
			sshFlags:               testSSHFlags,
			wrapperFlags:           testWrapperFlags,
			flags:                  testFlags,
			prepareTFClientMuxMock: defaultPrepareTFClientMuxMock,
			mockLoopStatusCheck: func(t *testing.T, expectedCtx context.Context, expectedTimeout time.Duration, expectedCheckForRunning bool) loopStatusCheckFn {
				return func(ctx context.Context, _ wrapper.LoopStatusCheckClient, timeout time.Duration, checkForRunning bool) error {
					assert.Equal(t, expectedCtx, ctx)
					assert.Equal(t, expectedTimeout, timeout)
					assert.Equal(t, expectedCheckForRunning, checkForRunning)

					return assert.AnError
				}
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, errWrapperStatusCheck)
			},
		},
		"tfClientMux loopStatusCheck success": {
			tfFlags:                testTFFlags,
			sshFlags:               testSSHFlags,
			wrapperFlags:           testWrapperFlags,
			flags:                  testFlags,
			prepareTFClientMuxMock: defaultPrepareTFClientMuxMock,
			mockLoopStatusCheck: func(t *testing.T, expectedCtx context.Context, expectedTimeout time.Duration, expectedCheckForRunning bool) loopStatusCheckFn {
				return func(ctx context.Context, _ wrapper.LoopStatusCheckClient, timeout time.Duration, checkForRunning bool) error {
					assert.Equal(t, expectedCtx, ctx)
					assert.Equal(t, expectedTimeout, timeout)
					assert.Equal(t, expectedCheckForRunning, checkForRunning)

					return nil
				}
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			require.NotNil(t, tt.prepareTFClientMuxMock, "prepareTFClientMuxMock must be defined in test definition")
			require.NotNil(t, tt.mockLoopStatusCheck, "mockLoopStatusCheck must be defined in test definition")

			tfc := newMockTfClientMux(t)
			tt.prepareTFClientMuxMock(t, tfc, ctx)

			s := New(logger.New(), nil, tt.tfFlags, tt.sshFlags, tt.wrapperFlags, tt.flags)
			s.tfClientMuxFactory = func(s *slog.Logger, client *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags wrapper.Flags) tfClientMux {
				assert.Equal(t, tt.tfFlags, tfFlags)
				assert.Equal(t, tt.sshFlags, sshFlags)
				assert.Equal(t, tt.wrapperFlags, wrapperFlags)

				return tfc
			}

			s.loopStatusCheck = tt.mockLoopStatusCheck(t, ctx, tt.flags.Timeout, expectedCheckForRunning)

			err := testFn(ctx, s)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestService_ExecuteWaitTerminated(t *testing.T) {
	testServiceExecute(t, wrapper.CheckForStopped, func(ctx context.Context, s *Service) error {
		return s.ExecuteWaitTerminated(ctx)
	})
}
