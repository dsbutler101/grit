package ssh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProxyJumpDialer(t *testing.T) {
	withSocket(t, func(t *testing.T, socketPath string) {
		withSSHServerCfg(t, func(t *testing.T, proxyCfg Config) {
			withSSHServerCfg(t, func(t *testing.T, targetCfg Config) {
				dialer, err := NewProxyJumpDialer(targetCfg, proxyCfg)
				require.NoError(t, err)

				runTestCall(t, dialer.Dial, socketPath)

				require.NoError(t, dialer.Close())
			})
		})
	})
}
