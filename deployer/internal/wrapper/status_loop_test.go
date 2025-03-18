package wrapper

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api"
)

func TestLoopStatusCheck(t *testing.T) {
	testTimeout := 50 * time.Millisecond

	tests := map[string]struct {
		prepareClientMock func(ctx context.Context, client *MockLoopStatusCheckClient)
		assertError       func(t *testing.T, err error)
	}{
		"gRPC request error failure": {
			prepareClientMock: func(ctx context.Context, client *MockLoopStatusCheckClient) {
				client.EXPECT().
					CheckStatus(ctx).
					Return(Status{Status: api.StatusUnknown}, assert.AnError)
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		"quick successful check": {
			prepareClientMock: func(ctx context.Context, client *MockLoopStatusCheckClient) {
				client.EXPECT().
					CheckStatus(ctx).
					Return(Status{Status: api.StatusRunning}, nil)
			},
		},
		"successful check after a retry": {
			prepareClientMock: func(ctx context.Context, client *MockLoopStatusCheckClient) {
				client.EXPECT().
					CheckStatus(ctx).
					Return(Status{Status: api.StatusUnknown}, nil).
					Once()

				client.EXPECT().
					CheckStatus(ctx).
					Return(Status{Status: api.StatusRunning}, nil).
					Once()
			},
		},
		"successful check after a retry when timeout potentially exceeded": {
			prepareClientMock: func(ctx context.Context, client *MockLoopStatusCheckClient) {
				client.EXPECT().
					CheckStatus(ctx).
					Return(Status{Status: api.StatusUnknown}, nil).
					Once()

				client.EXPECT().
					CheckStatus(ctx).
					Return(Status{Status: api.StatusRunning}, nil).
					Run(func(_ context.Context) {
						time.Sleep(2 * testTimeout)
					}).
					Once()
			},
		},
		"unsuccessful check after a retry when timeout exceeded": {
			prepareClientMock: func(ctx context.Context, client *MockLoopStatusCheckClient) {
				client.EXPECT().
					CheckStatus(ctx).
					Return(Status{Status: api.StatusUnknown}, nil).
					Once()

				client.EXPECT().
					CheckStatus(ctx).
					Return(Status{Status: api.StatusUnknown}, nil).
					Run(func(_ context.Context) {
						time.Sleep(2 * testTimeout)
					}).
					Once()
			},
			assertError: func(t *testing.T, err error) {
				var eerr *StatusCheckLoopTimeoutExceededError
				if assert.ErrorAs(t, err, &eerr) {
					assert.Equal(t, testTimeout, eerr.timeout)
				}
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			clientMock := NewMockLoopStatusCheckClient(t)

			require.NotNil(t, tt.prepareClientMock, "prepareClientMock function in the test definition must be defined")
			tt.prepareClientMock(ctx, clientMock)

			err := loopStatusCheckWithSleep(ctx, clientMock, testTimeout, CheckForRunning, 10*time.Millisecond)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
