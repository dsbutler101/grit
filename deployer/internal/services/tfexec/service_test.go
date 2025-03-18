package tfexec

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/logger"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
)

func TestService_ExecuteUp(t *testing.T) {
	testTarget := "test-target"

	tests := map[string]struct {
		prepareTFClientMock func(t *testing.T, m *mockTfClient, expectedCtx context.Context)
		assertError         func(t *testing.T, err error)
	}{
		"init failure": {
			prepareTFClientMock: func(t *testing.T, m *mockTfClient, expectedCtx context.Context) {
				m.EXPECT().Init(expectedCtx, testTarget).Return(assert.AnError)
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, errInitializingTerraform)
			},
		},
		"destroy failure": {
			prepareTFClientMock: func(t *testing.T, m *mockTfClient, expectedCtx context.Context) {
				m.EXPECT().Init(expectedCtx, testTarget).Return(nil)
				m.EXPECT().Apply(expectedCtx, testTarget).Return(assert.AnError)
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, errApplyingTerraformResources)
			},
		},
		"success": {
			prepareTFClientMock: func(t *testing.T, m *mockTfClient, expectedCtx context.Context) {
				m.EXPECT().Init(expectedCtx, testTarget).Return(nil)
				m.EXPECT().Apply(expectedCtx, testTarget).Return(nil)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			require.NotNil(t, tt.prepareTFClientMock, "prepareTFClientMock must be defined in test definition")

			tfc := newMockTfClient(t)
			tt.prepareTFClientMock(t, tfc, ctx)

			s := New(logger.New(), tfc, terraform.Flags{
				Target: testTarget,
			})

			err := s.ExecuteUp(ctx)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestService_ExecuteDown(t *testing.T) {
	testTarget := "test-target"

	tests := map[string]struct {
		prepareTFClientMock func(t *testing.T, m *mockTfClient, expectedCtx context.Context)
		assertError         func(t *testing.T, err error)
	}{
		"init failure": {
			prepareTFClientMock: func(t *testing.T, m *mockTfClient, expectedCtx context.Context) {
				m.EXPECT().Init(expectedCtx, testTarget).Return(assert.AnError)
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, errInitializingTerraform)
			},
		},
		"destroy failure": {
			prepareTFClientMock: func(t *testing.T, m *mockTfClient, expectedCtx context.Context) {
				m.EXPECT().Init(expectedCtx, testTarget).Return(nil)
				m.EXPECT().Destroy(expectedCtx, testTarget).Return(assert.AnError)
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, errDestroyingTerraformResources)
			},
		},
		"success": {
			prepareTFClientMock: func(t *testing.T, m *mockTfClient, expectedCtx context.Context) {
				m.EXPECT().Init(expectedCtx, testTarget).Return(nil)
				m.EXPECT().Destroy(expectedCtx, testTarget).Return(nil)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			require.NotNil(t, tt.prepareTFClientMock, "prepareTFClientMock must be defined in test definition")

			tfc := newMockTfClient(t)
			tt.prepareTFClientMock(t, tfc, ctx)

			s := New(logger.New(), tfc, terraform.Flags{
				Target: testTarget,
			})

			err := s.ExecuteDown(ctx)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
