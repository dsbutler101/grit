package ssh

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandDialer(t *testing.T) {
	withSocket(t, func(t *testing.T, socketPath string) {
		withCompiledNC(t, func(t *testing.T, ncPath string) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			dialer := NewCommandDialer(CommandDef{
				Cmd:  ncPath,
				Args: []string{"--unix-socket", socketPath},
			})

			err := dialer.Start(ctx)
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
}
