package wrapper

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/logger"
)

func TestClient_Connect(t *testing.T) {
	testTimeout := 10 * time.Millisecond
	testDialer := func(network string, address string) (net.Conn, error) {
		return net.Dial(network, address)
	}

	assertGRPCConnectionRetryExceededError := func(t *testing.T, err error) {
		var eerr *GRPCConnectionRetryExceededError
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

			c, err := NewClient(logger.New(), testDialer, "unix:///tmp/test.sock")
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
