package terraform

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultCommander_run(t *testing.T) {
	testExitCode := 123

	tests := map[string]struct {
		args         []string
		assertError  func(t *testing.T, err error)
		assertStdout func(t *testing.T, out string)
		assertStderr func(t *testing.T, out string)
	}{
		"command failure": {
			args: []string{"fail", strconv.Itoa(testExitCode)},
			assertError: func(t *testing.T, err error) {
				var eerr *exec.ExitError
				if assert.ErrorAs(t, err, &eerr) {
					assert.Equal(t, testExitCode, eerr.ProcessState.ExitCode())
				}
			},
			assertStderr: func(t *testing.T, out string) {
				assert.Contains(t, out, fmt.Sprintf("Exiting with %d", testExitCode))
			},
		},
		"command success": {
			assertStdout: func(t *testing.T, out string) {
				assert.Equal(t, "out\n", out)
			},
			assertStderr: func(t *testing.T, out string) {
				assert.Equal(t, "err\n", out)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			outBuf := bytes.NewBuffer(nil)
			errBuf := bytes.NewBuffer(nil)

			c := newDefaultCommander(outBuf, errBuf, t.TempDir())
			err := withCompiledTestCommanderBinary(t, func(binaryPath string) error {
				return c.run(ctx, binaryPath, tt.args...)
			})

			if tt.assertStdout != nil {
				tt.assertStdout(t, outBuf.String())
			}

			if tt.assertStderr != nil {
				tt.assertStderr(t, errBuf.String())
			}

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func withCompiledTestCommanderBinary(t *testing.T, fn func(binaryPath string) error) error {
	testBinary := filepath.Join(t.TempDir(), "cmd.exe")
	require.NoError(t, exec.Command("go", "build", "-o", testBinary, "./testdata/commander/").Run())

	return fn(testBinary)
}

func TestDefaultCommander_exitCode(t *testing.T) {
	testExitCode := 123

	tests := map[string]struct {
		runCommand       bool
		args             []string
		expectedExitCode int
	}{
		"command not executed": {
			runCommand:       false,
			expectedExitCode: CommanderCommandNotExecutedExitCode,
		},
		"command succeeded": {
			runCommand:       true,
			expectedExitCode: 0,
		},
		"command failed": {
			runCommand:       true,
			args:             []string{"fail", strconv.Itoa(testExitCode)},
			expectedExitCode: testExitCode,
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			c := newDefaultCommander(io.Discard, io.Discard, t.TempDir())
			_ = withCompiledTestCommanderBinary(t, func(binaryPath string) error {
				if !tt.runCommand {
					return nil
				}

				return c.run(ctx, binaryPath, tt.args...)
			})

			assert.Equal(t, tt.expectedExitCode, c.exitCode())
		})
	}
}
