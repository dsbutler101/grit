package ssh

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProxyCommandDialer(t *testing.T) {
	withSocket(t, func(t *testing.T, socketPath string) {
		withCompiledNC(t, func(t *testing.T, ncPath string) {
			withSSHServerCfg(t, func(t *testing.T, cfg Config) {
				ctx, cancelFn := context.WithCancel(context.Background())
				defer cancelFn()

				dialer, err := NewProxyCommandDialer(cfg, CommandDef{
					Cmd:  ncPath,
					Args: []string{"--tcp-socket", cfg.Address},
				})
				require.NoError(t, err)

				err = dialer.Start(ctx)
				require.NoError(t, err)

				dialerErrChan := make(chan error)
				go func() {
					dialerErrChan <- dialer.Wait()
				}()

				runTestCall(t, dialer.Dial, socketPath)

				err = dialer.Close()
				require.NoError(t, err)

				require.NoError(t, <-dialerErrChan)
			})
		})
	})
}
