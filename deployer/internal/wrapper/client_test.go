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
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	dialer := func(network string, address string) (net.Conn, error) {
		return net.Dial(network, address)
	}

	c, err := NewClient(logger.New(), dialer, "unix:///tmp/test.sock")
	require.NoError(t, err)

	err = c.Connect(ctx, 1*time.Second)
	require.NoError(t, err)

	status, err := c.CheckStatus(ctx)
	assert.NoError(t, err)
	t.Log(status)
	assert.NoError(t, c.InitGracefulShutdown(ctx))
}
