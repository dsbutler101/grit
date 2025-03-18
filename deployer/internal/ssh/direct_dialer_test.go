package ssh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirectDialer(t *testing.T) {
	withSocket(t, func(t *testing.T, socketPath string) {
		withSSHServerCfg(t, func(t *testing.T, cfg Config) {
			dialer, err := NewDirectDialer(cfg)
			require.NoError(t, err)

			runTestCall(t, dialer.Dial, socketPath)

			require.NoError(t, dialer.Close())
		})
	})
}
