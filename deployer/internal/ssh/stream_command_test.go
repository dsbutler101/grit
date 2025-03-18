package ssh

import (
	"context"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultExecCmd_Kill(t *testing.T) {
	withCompiledNC(t, func(t *testing.T, ncPath string) {
		testSocket := filepath.Join(t.TempDir(), "test.sock")

		tests := map[string]struct {
			mockCmd func(t *testing.T, binaryPath string) *exec.Cmd
		}{
			"cmd not set": {
				mockCmd: func(_ *testing.T, binaryPath string) *exec.Cmd {
					return nil
				},
			},
			"cmd.Process is nil": {
				mockCmd: func(t *testing.T, binaryPath string) *exec.Cmd {
					return exec.Command("unknown")
				},
			},
			"cmd.Process is not nil": {
				mockCmd: func(t *testing.T, binaryPath string) *exec.Cmd {
					cmd := exec.Command(binaryPath, "--unix-socket", testSocket)
					require.NoError(t, cmd.Start())
					time.Sleep(100 * time.Millisecond)

					return cmd
				},
			},
		}

		for tn, tt := range tests {
			t.Run(tn, func(t *testing.T) {
				require.NotNil(t, tt.mockCmd, "mockCmd must be defined in test definition")
				dec := defaultExecCmd{Cmd: tt.mockCmd(t, ncPath)}
				assert.NoError(t, dec.Kill())
			})
		}
	})
}

func TestStreamCommand_doubleStartCall(t *testing.T) {
	withCompiledNC(t, func(t *testing.T, ncPath string) {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()

		socketPath := filepath.Join(t.TempDir(), "test.sock")
		c := newStreamCommand(ncPath, "--unix-socket", socketPath)

		assert.NoError(t, c.Start(ctx))
		time.Sleep(100 * time.Millisecond)
		assert.ErrorIs(t, c.Start(ctx), errStreamCommandAlreadyRunning)
		assert.NoError(t, c.Kill())
	})
}

func TestStreamCommand_Start(t *testing.T) {
	testCommand := "test-cmd"
	testArgs := []string{"test-arg-1", "test-arg-2"}

	tests := map[string]struct {
		prepareExecCmdMock func(t *testing.T, m *mockExecCmd)
		assertError        func(t *testing.T, err error)
	}{
		"StdoutPipe() failure": {
			prepareExecCmdMock: func(t *testing.T, m *mockExecCmd) {
				m.EXPECT().StdoutPipe().Return(nil, assert.AnError)
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errStreamCommandCreatingStdoutPipe)
			},
		},
		"StderrPipe() failure": {
			prepareExecCmdMock: func(t *testing.T, m *mockExecCmd) {
				m.EXPECT().StdoutPipe().Return(nil, nil)
				m.EXPECT().StderrPipe().Return(nil, assert.AnError)
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errStreamCommandCreatingStderrPipe)
			},
		},
		"StdinPipe() failure": {
			prepareExecCmdMock: func(t *testing.T, m *mockExecCmd) {
				m.EXPECT().StdoutPipe().Return(nil, nil)
				m.EXPECT().StderrPipe().Return(nil, nil)
				m.EXPECT().StdinPipe().Return(nil, assert.AnError)
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errStreamCommandCreatingStdinPipe)
			},
		},
		"cmd.Start() failure": {
			prepareExecCmdMock: func(t *testing.T, m *mockExecCmd) {
				m.EXPECT().StdoutPipe().Return(nil, nil)
				m.EXPECT().StderrPipe().Return(nil, nil)
				m.EXPECT().StdinPipe().Return(nil, nil)
				m.EXPECT().Start().Return(assert.AnError)
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		"cmd.Start() success": {
			prepareExecCmdMock: func(t *testing.T, m *mockExecCmd) {
				m.EXPECT().StdoutPipe().Return(nil, nil)
				m.EXPECT().StderrPipe().Return(nil, nil)
				m.EXPECT().StdinPipe().Return(nil, nil)
				m.EXPECT().Start().Return(nil)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			ecMock := newMockExecCmd(t)

			require.NotNil(t, tt.prepareExecCmdMock, "prepareExecCmdMock must be defined in test definition")
			tt.prepareExecCmdMock(t, ecMock)

			sc := newStreamCommand(testCommand, testArgs...)
			sc.commandFactory = func(cctx context.Context, command string, args ...string) execCmd {
				assert.Equal(t, testCommand, command)
				assert.Equal(t, testArgs, args)
				assert.Equal(t, ctx, cctx)

				return ecMock
			}

			err := sc.Start(ctx)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestStreamCommand_Connection(t *testing.T) {
	testCommand := "test-cmd"
	testArgs := []string{"test-arg-1", "test-arg-2"}

	tests := map[string]struct {
		start            bool
		assertConnection func(t *testing.T, c net.Conn)
	}{
		"command not started": {
			start: false,
			assertConnection: func(t *testing.T, c net.Conn) {
				assert.Nil(t, c)
			},
		},
		"command started": {
			start: true,
			assertConnection: func(t *testing.T, c net.Conn) {
				assert.NotNil(t, c)
				assert.IsType(t, &streamConn{}, c)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			ecMock := newMockExecCmd(t)

			sc := newStreamCommand(testCommand, testArgs...)
			sc.commandFactory = func(cctx context.Context, command string, args ...string) execCmd {
				assert.Equal(t, testCommand, command)
				assert.Equal(t, testArgs, args)
				assert.Equal(t, ctx, cctx)

				return ecMock
			}

			if tt.start {
				ecMock.EXPECT().StdoutPipe().Return(nil, nil)
				ecMock.EXPECT().StderrPipe().Return(nil, nil)
				ecMock.EXPECT().StdinPipe().Return(nil, nil)
				ecMock.EXPECT().Start().Return(nil)
				require.NoError(t, sc.Start(ctx))
			}

			require.NotNil(t, tt.assertConnection, "assertConnection must be defined in test definition")
			tt.assertConnection(t, sc.Connection())
		})
	}
}

func TestStreamCommand_Wait(t *testing.T) {
	testCommand := "test-cmd"
	testArgs := []string{"test-arg-1", "test-arg-2"}

	tests := map[string]struct {
		start       bool
		waitErr     error
		assertError func(t *testing.T, err error)
	}{
		"cmd not started": {
			start: false,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errStreamCommandNotStarted)
			},
		},
		"cmd.Wait() failure": {
			start:   true,
			waitErr: assert.AnError,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		"cmd.Wait() exec.ExitError failure": {
			start: true,
			waitErr: &exec.ExitError{
				ProcessState: &os.ProcessState{},
			},
			assertError: func(t *testing.T, err error) {
				var eerr *exec.ExitError
				if assert.ErrorAs(t, err, &eerr) {
					assert.Equal(t, 0, eerr.ExitCode())
				}
			},
		},
		"cmd.Wait() success": {
			start:   true,
			waitErr: nil,
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			ecMock := newMockExecCmd(t)

			sc := newStreamCommand(testCommand, testArgs...)
			sc.commandFactory = func(cctx context.Context, command string, args ...string) execCmd {
				assert.Equal(t, testCommand, command)
				assert.Equal(t, testArgs, args)
				assert.Equal(t, ctx, cctx)

				return ecMock
			}

			if tt.start {
				ecMock.EXPECT().StdoutPipe().Return(nil, nil)
				ecMock.EXPECT().StderrPipe().Return(nil, nil)
				ecMock.EXPECT().StdinPipe().Return(nil, nil)
				ecMock.EXPECT().Start().Return(nil)
				require.NoError(t, sc.Start(ctx))

				ecMock.EXPECT().Wait().Return(tt.waitErr)
			}

			err := sc.Wait()

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestStreamCommand_WaitWithCmdFailures(t *testing.T) {
	testExitCode := 100
	withCompiledExitCode(t, func(t *testing.T, exitCodePath string) {
		tests := map[string]struct {
			commandArgs []string
			kill        bool
			assertError func(t *testing.T, err error)
		}{
			"command succeeded": {},
			"command failed": {
				commandArgs: []string{strconv.Itoa(testExitCode)},
				assertError: func(t *testing.T, err error) {
					var eerr *exec.ExitError
					if assert.ErrorAs(t, err, &eerr) {
						assert.Equal(t, testExitCode, eerr.ExitCode())
					}
				},
			},
			"command killed": {
				kill: true,
			},
		}

		for tn, tt := range tests {
			t.Run(tn, func(t *testing.T) {
				ctx, cancelFn := context.WithCancel(context.Background())
				defer cancelFn()

				var err error

				sc := newStreamCommand(exitCodePath, tt.commandArgs...)
				require.NoError(t, sc.Start(ctx))

				if tt.kill {
					assert.NoError(t, sc.Kill())
				}

				err = sc.Wait()

				if tt.assertError != nil {
					tt.assertError(t, err)
					return
				}

				assert.NoError(t, err)
			})
		}
	})
}

func withCompiledExitCode(t *testing.T, fn func(t *testing.T, exitCodePath string)) {
	sourcePath := "." + string(filepath.Separator) + filepath.Join("testdata", "exitcode")
	binPath := filepath.Join(t.TempDir(), "ec.exe")

	cmd := exec.Command("go", "build", "-o", binPath, sourcePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	require.NoError(t, err, "Compiling exitcode binary")

	fn(t, binPath)
}
