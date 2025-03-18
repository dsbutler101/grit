package wrapper

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/logger"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api/client"
)

var (
	testNOOPDialer = func(network string, address string) (net.Conn, error) {
		return nil, nil
	}
)

func TestClient_Connect(t *testing.T) {
	testTimeout := 10 * time.Millisecond

	assertGRPCConnectionRetryExceededError := func(t *testing.T, err error) {
		var eerr *GRPCConnectionWaitTimeoutExceededError
		if assert.ErrorAs(t, err, &eerr) {
			assert.Equal(t, testTimeout, eerr.timeout)
			assert.ErrorIs(t, eerr.err, assert.AnError)
		}
	}

	tests := map[string]struct {
		connectError error
		assertError  func(t *testing.T, err error)
	}{
		"connection succeeds": {
			connectError: nil,
			assertError:  nil,
		},
		"connection error": {
			connectError: assert.AnError,
			assertError:  assertGRPCConnectionRetryExceededError,
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			grpcClientMock := newMockGrpcClient(t)
			grpcClientMock.EXPECT().
				ConnectWithTimeout(ctx, testTimeout).
				Return(tt.connectError)

			c, err := NewClient(logger.New(), testNOOPDialer, "unix:///tmp/test.sock")
			require.NoError(t, err)

			c.c = grpcClientMock

			err = c.Connect(ctx, testTimeout)
			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestClient_CheckStatus(t *testing.T) {
	testFailureReason := "test-failure-reason"

	tests := map[string]struct {
		grpcResponse client.CheckStatusResponse
		requestError error
		assertError  func(t *testing.T, err error)
		assertStatus func(t *testing.T, status Status)
	}{
		"status check successful": {
			grpcResponse: client.CheckStatusResponse{
				Status:        api.StatusRunning,
				FailureReason: testFailureReason,
			},
			requestError: nil,
			assertError:  nil,
			assertStatus: func(t *testing.T, status Status) {
				assert.Equal(t, api.StatusRunning, status.Status)
				assert.Equal(t, testFailureReason, status.FailureReason)
			},
		},
		"check status request failure": {
			requestError: assert.AnError,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrCheckStatus)
				assert.Contains(t, err.Error(), assert.AnError.Error())
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			grpcClientMock := newMockGrpcClient(t)
			grpcClientMock.EXPECT().CheckStatus(ctx).Return(tt.grpcResponse, tt.requestError)

			c, err := NewClient(logger.New(), testNOOPDialer, "unix:///tmp/test.sock")
			require.NoError(t, err)

			c.c = grpcClientMock

			status, err := c.CheckStatus(ctx)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)

			if tt.assertStatus != nil {
				tt.assertStatus(t, status)
			}
		})
	}
}

func TestClient_InitGracefulShutdown(t *testing.T) {
	tests := map[string]struct {
		gRPCResponse client.CheckStatusResponse
		requestError error
		assertError  func(t *testing.T, err error)
	}{
		"graceful shutdown initialization successful": {
			gRPCResponse: client.CheckStatusResponse{
				Status:        api.StatusRunning,
				FailureReason: "",
			},
			requestError: nil,
			assertError:  nil,
		},
		"graceful shutdown initialization failure": {
			requestError: assert.AnError,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrInitGracefulShutdown)
				assert.Contains(t, err.Error(), assert.AnError.Error())
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			grpcClientMock := newMockGrpcClient(t)
			grpcClientMock.EXPECT().InitGracefulShutdown(ctx, mock.Anything).Return(tt.gRPCResponse, tt.requestError)

			c, err := NewClient(logger.New(), testNOOPDialer, "unix:///tmp/test.sock")
			require.NoError(t, err)

			c.c = grpcClientMock

			err = c.InitGracefulShutdown(ctx)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestClient_InitForcefulShutdown(t *testing.T) {
	tests := map[string]struct {
		gRPCResponse client.CheckStatusResponse
		requestError error
		assertError  func(t *testing.T, err error)
	}{
		"graceful shutdown initialization successful": {
			gRPCResponse: client.CheckStatusResponse{
				Status:        api.StatusRunning,
				FailureReason: "",
			},
			requestError: nil,
			assertError:  nil,
		},
		"graceful shutdown initialization failure": {
			requestError: assert.AnError,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrInitForcefulShutdown)
				assert.Contains(t, err.Error(), assert.AnError.Error())
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			grpcClientMock := newMockGrpcClient(t)
			grpcClientMock.EXPECT().InitForcefulShutdown(ctx).Return(tt.gRPCResponse, tt.requestError)

			c, err := NewClient(logger.New(), testNOOPDialer, "unix:///tmp/test.sock")
			require.NoError(t, err)

			c.c = grpcClientMock

			err = c.InitForcefulShutdown(ctx)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
