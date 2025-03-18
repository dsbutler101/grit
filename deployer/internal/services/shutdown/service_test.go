package shutdown

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

func TestService_Execute(t *testing.T) {
	testTimeout := 1 * time.Millisecond
	testTFFlags := terraform.Flags{
		Target: "test-target",
	}
	testSSHFlags := ssh.Flags{
		Username: "test-username",
	}
	testWrapperFlags := wrapper.Flags{
		ConnectionTimeout: testTimeout,
	}
	testFlags := Flags{
		Forceful: false,
	}

	defaultPrepareTFClientMuxMock := func(t *testing.T, m *mockTfClientMux, expectedCtx context.Context, cliMock *wrapper.MockCallbackClient) {
		m.EXPECT().Execute(expectedCtx, mock.Anything).RunAndReturn(func(ctx context.Context, fn wrapper.CallbackFn) error {
			return fn(ctx, cliMock)
		})
	}

	tests := map[string]struct {
		tfFlags                   terraform.Flags
		sshFlags                  ssh.Flags
		wrapperFlags              wrapper.Flags
		flags                     Flags
		prepareCallbackClientMock func(t *testing.T, m *wrapper.MockCallbackClient, expectedCtx context.Context)
		prepareTFClientMuxMock    func(t *testing.T, m *mockTfClientMux, expectedCtx context.Context, cliMock *wrapper.MockCallbackClient)
		assertError               func(t *testing.T, err error)
	}{
		"tfClientMux execution error": {
			tfFlags:                   testTFFlags,
			sshFlags:                  testSSHFlags,
			wrapperFlags:              testWrapperFlags,
			flags:                     testFlags,
			prepareCallbackClientMock: func(t *testing.T, m *wrapper.MockCallbackClient, expectedCtx context.Context) {},
			prepareTFClientMuxMock: func(t *testing.T, m *mockTfClientMux, expectedCtx context.Context, cliMock *wrapper.MockCallbackClient) {
				m.EXPECT().Execute(expectedCtx, mock.Anything).Return(assert.AnError)
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		"forceful shutdown error": {
			tfFlags:      testTFFlags,
			sshFlags:     testSSHFlags,
			wrapperFlags: testWrapperFlags,
			flags: Flags{
				Forceful: true,
			},
			prepareCallbackClientMock: func(t *testing.T, m *wrapper.MockCallbackClient, expectedCtx context.Context) {
				m.EXPECT().InitForcefulShutdown(expectedCtx).Return(assert.AnError)
			},
			prepareTFClientMuxMock: defaultPrepareTFClientMuxMock,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, errInitiatingForcefulShutdown)
			},
		},
		"forceful shutdown success": {
			tfFlags:      testTFFlags,
			sshFlags:     testSSHFlags,
			wrapperFlags: testWrapperFlags,
			flags: Flags{
				Forceful: true,
			},
			prepareCallbackClientMock: func(t *testing.T, m *wrapper.MockCallbackClient, expectedCtx context.Context) {
				m.EXPECT().InitForcefulShutdown(expectedCtx).Return(nil)
			},
			prepareTFClientMuxMock: defaultPrepareTFClientMuxMock,
		},
		"graceful shutdown error": {
			tfFlags:      testTFFlags,
			sshFlags:     testSSHFlags,
			wrapperFlags: testWrapperFlags,
			flags:        testFlags,
			prepareCallbackClientMock: func(t *testing.T, m *wrapper.MockCallbackClient, expectedCtx context.Context) {
				m.EXPECT().InitGracefulShutdown(expectedCtx).Return(assert.AnError)
			},
			prepareTFClientMuxMock: defaultPrepareTFClientMuxMock,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, errInitiatingGracefulShutdown)
			},
		},
		"graceful shutdown success": {
			tfFlags:      testTFFlags,
			sshFlags:     testSSHFlags,
			wrapperFlags: testWrapperFlags,
			flags:        testFlags,
			prepareCallbackClientMock: func(t *testing.T, m *wrapper.MockCallbackClient, expectedCtx context.Context) {
				m.EXPECT().InitGracefulShutdown(expectedCtx).Return(nil)
			},
			prepareTFClientMuxMock: defaultPrepareTFClientMuxMock,
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithTimeout(context.Background(), testTimeout)
			defer cancelFn()

			require.NotNil(t, tt.prepareTFClientMuxMock, "prepareTFClientMuxMock must be defined in test definition")
			require.NotNil(t, tt.prepareCallbackClientMock, "prepareCallbackClientMock must be defined in test definition")

			cc := wrapper.NewMockCallbackClient(t)
			tt.prepareCallbackClientMock(t, cc, ctx)

			tfc := newMockTfClientMux(t)
			tt.prepareTFClientMuxMock(t, tfc, ctx, cc)

			s := New(logger.New(), nil, tt.tfFlags, tt.sshFlags, tt.wrapperFlags, tt.flags)
			s.tfClientMuxFactory = func(s *slog.Logger, client *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags wrapper.Flags) tfClientMux {
				assert.Equal(t, tt.tfFlags, tfFlags)
				assert.Equal(t, tt.sshFlags, sshFlags)
				assert.Equal(t, tt.wrapperFlags, wrapperFlags)

				return tfc
			}

			err := s.Execute(ctx)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
