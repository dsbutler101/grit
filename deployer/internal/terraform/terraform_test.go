package terraform

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/logger"
)

func TestCommandError(t *testing.T) {
	testCommand := "show"
	testExitCode := 1

	err := NewCommandError(testCommand, testExitCode, assert.AnError)
	assert.Equal(t, testExitCode, err.ExitCode())
	assert.Equal(t, assert.AnError, err.Unwrap())
}

func TestClient_Init(t *testing.T) {
	testRunTerraformWithOutput(t, func(t *testing.T) testTerraformCommandRunDef {
		workdir := t.TempDir()

		return testTerraformCommandRunDef{
			tfCommand: tfCommandInit,
			fn: func(ctx context.Context, c *Client) error {
				return c.Init(ctx, workdir)
			},
		}
	})
}

func TestClient_Apply(t *testing.T) {
	testRunTerraformWithOutput(t, func(t *testing.T) testTerraformCommandRunDef {
		workdir := t.TempDir()

		return testTerraformCommandRunDef{
			tfCommand: tfCommandApply,
			tfArgs:    []interface{}{tfAutoApproveFlag},
			fn: func(ctx context.Context, c *Client) error {
				return c.Apply(ctx, workdir)
			},
		}
	})
}

func TestClient_Destroy(t *testing.T) {
	testRunTerraformWithOutput(t, func(t *testing.T) testTerraformCommandRunDef {
		workdir := t.TempDir()

		return testTerraformCommandRunDef{
			tfCommand: tfCommandDestroy,
			tfArgs:    []interface{}{tfAutoApproveFlag},
			fn: func(ctx context.Context, c *Client) error {
				return c.Destroy(ctx, workdir)
			},
		}
	})
}

type testTerraformCommandRunDef struct {
	tfCommand string
	tfArgs    []interface{}
	fn        func(ctx context.Context, c *Client) error
}

func testRunTerraformWithOutput(t *testing.T, fn func(t *testing.T) testTerraformCommandRunDef) {
	testExitCode := 5
	testExecPath := filepath.Join(t.TempDir(), "terraform")

	tests := map[string]struct {
		execPath       string
		prepareMockRun func(t *testing.T, cmdr *mockCommander, expectedContext context.Context, def testTerraformCommandRunDef)
		assertError    func(t *testing.T, err error, def testTerraformCommandRunDef)
	}{
		"exec path not set": {
			assertError: func(t *testing.T, err error, _ testTerraformCommandRunDef) {
				assert.ErrorIs(t, err, ErrTerraformExecPathNotSet)
			},
		},
		"terraform command execution failure": {
			execPath: testExecPath,
			prepareMockRun: func(t *testing.T, cmdr *mockCommander, expectedContext context.Context, def testTerraformCommandRunDef) {
				args := append([]interface{}{def.tfCommand}, def.tfArgs...)

				cmdr.EXPECT().run(expectedContext, testExecPath, args...).Return(assert.AnError)
				cmdr.EXPECT().exitCode().Return(testExitCode)
			},
			assertError: func(t *testing.T, err error, def testTerraformCommandRunDef) {
				var eerr *CommandError
				if assert.ErrorAs(t, err, &eerr) {
					assert.Equal(t, testExitCode, eerr.ExitCode())
					assert.Equal(t, def.tfCommand, eerr.command)
				}
			},
		},
		"terraform command execution succeeded": {
			execPath: testExecPath,
			prepareMockRun: func(t *testing.T, cmdr *mockCommander, expectedContext context.Context, def testTerraformCommandRunDef) {
				args := append([]interface{}{def.tfCommand}, def.tfArgs...)

				cmdr.EXPECT().run(expectedContext, testExecPath, args...).Return(nil)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			require.NotNil(t, fn)
			def := fn(t)

			cmdr := newMockCommander(t)

			c := New(logger.New())
			c.SetExecPath(tt.execPath)

			c.commanderFactory = func(stdout io.Writer, stderr io.Writer, workdir string) commander {
				return cmdr
			}

			if tt.prepareMockRun != nil {
				tt.prepareMockRun(t, cmdr, ctx, def)
			}

			require.NotNil(t, def.fn, "def.fn must be defined")

			err := def.fn(ctx, c)
			if tt.assertError != nil {
				tt.assertError(t, err, def)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestClient_ReadStateDir(t *testing.T) {
	testExitCode := 3
	testExecPath := filepath.Join(t.TempDir(), "terraform")
	testStateOutput := "test-state-output"
	testRM1Alias := "test-runner-manager-1"
	testInstanceName := "test-runner-manager-1-instance"
	testRMs := RunnerManagers{
		testRM1Alias: RunnerManager{
			InstanceName: testInstanceName,
		},
	}

	noopMockStateReader := func(t *testing.T) stateReader {
		return func(buf *bytes.Buffer) (RunnerManagers, error) {
			return nil, nil
		}
	}

	defaultPrepareCommanderMock := func(t *testing.T, cmdr *mockCommander, expectedContext context.Context) {
		var out io.Writer
		cmdr.EXPECT().setStdout(mock.Anything).Run(func(stdout io.Writer) {
			out = stdout
		})
		cmdr.EXPECT().
			run(expectedContext, testExecPath, tfCommandInit).
			Return(nil)
		cmdr.EXPECT().
			run(expectedContext, testExecPath, tfCommandShow, "-json").
			Run(func(_ context.Context, _ string, _ ...string) {
				_, err := fmt.Fprint(out, testStateOutput)
				assert.NoError(t, err)
			}).
			Return(nil)
	}

	tests := map[string]struct {
		prepareCommanderMock func(t *testing.T, cmdr *mockCommander, expectedContext context.Context)
		mockStateReader      func(t *testing.T) stateReader
		assertError          func(t *testing.T, err error)
		assertRMs            func(t *testing.T, rms RunnerManagers)
	}{
		"reading state from dir fails": {
			prepareCommanderMock: func(t *testing.T, cmdr *mockCommander, expectedContext context.Context) {
				cmdr.EXPECT().setStdout(mock.Anything)
				cmdr.EXPECT().
					run(expectedContext, testExecPath, tfCommandInit).
					Return(nil)
				cmdr.EXPECT().
					run(expectedContext, testExecPath, tfCommandShow, "-json").
					Return(assert.AnError)
				cmdr.EXPECT().exitCode().Return(testExitCode)
			},
			mockStateReader: noopMockStateReader,
			assertError: func(t *testing.T, err error) {
				var eerr *CommandError
				if assert.ErrorAs(t, err, &eerr) {
					assert.Equal(t, testExitCode, eerr.ExitCode())
					assert.ErrorIs(t, eerr, assert.AnError)
				}
			},
		},
		"reading state fails": {
			prepareCommanderMock: defaultPrepareCommanderMock,
			mockStateReader: func(t *testing.T) stateReader {
				return func(buf *bytes.Buffer) (RunnerManagers, error) {
					assert.Equal(t, testStateOutput, buf.String())

					return nil, assert.AnError
				}
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		"reading state succeeds": {
			prepareCommanderMock: defaultPrepareCommanderMock,
			mockStateReader: func(t *testing.T) stateReader {
				return func(buf *bytes.Buffer) (RunnerManagers, error) {
					assert.Equal(t, testStateOutput, buf.String())

					return testRMs, nil
				}
			},
			assertRMs: func(t *testing.T, rms RunnerManagers) {
				assert.Equal(t, testRMs, rms)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			cmdr := newMockCommander(t)

			require.NotNil(t, tt.prepareCommanderMock, "prepareCommanderMock must be defined in test definition")
			tt.prepareCommanderMock(t, cmdr, ctx)

			c := New(logger.New())
			c.SetExecPath(testExecPath)
			c.commanderFactory = func(stdout io.Writer, stderr io.Writer, workdir string) commander {
				cmdr.setStdout(stdout)

				return cmdr
			}

			require.NotNil(t, tt.mockStateReader, "mockStateReader must be defined in test definition")
			c.stateReader = tt.mockStateReader(t)

			rms, err := c.ReadStateDir(ctx, t.TempDir())

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, tt.assertRMs, "assertRMs must be defined in test definition")
			tt.assertRMs(t, rms)
		})
	}
}

func TestClient_ReadStateFile(t *testing.T) {
	testStateFilePath := filepath.Join(t.TempDir(), "terraform-state-file")
	testUnknownStateFile := filepath.Join(t.TempDir(), "terraform-unknown-state-file")
	testRM1Alias := "test-runner-manager-1"
	testInstanceName := "test-runner-manager-1-instance"
	testRMs := RunnerManagers{
		testRM1Alias: RunnerManager{
			InstanceName: testInstanceName,
		},
	}

	testStateOutput := "test-state-output"
	err := os.WriteFile(testStateFilePath, []byte(testStateOutput), 0600)
	require.NoError(t, err)

	noopMockStateReader := func(t *testing.T) stateReader {
		return func(buf *bytes.Buffer) (RunnerManagers, error) {
			return nil, nil
		}
	}

	tests := map[string]struct {
		stateFilePath   string
		mockStateReader func(t *testing.T) stateReader
		assertError     func(t *testing.T, err error)
		assertRMs       func(t *testing.T, rms RunnerManagers)
	}{
		"file reading error": {
			stateFilePath:   testUnknownStateFile,
			mockStateReader: noopMockStateReader,
			assertError: func(t *testing.T, err error) {
				var eerr *os.PathError
				if assert.ErrorAs(t, err, &eerr) {
					assert.Equal(t, testUnknownStateFile, eerr.Path)
				}

				assert.ErrorIs(t, err, errReadingTerraformStateFile)
			},
		},
		"reading state fails": {
			stateFilePath: testStateFilePath,
			mockStateReader: func(t *testing.T) stateReader {
				return func(buf *bytes.Buffer) (RunnerManagers, error) {
					assert.Equal(t, testStateOutput, buf.String())

					return nil, assert.AnError
				}
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		"reading state succeeded": {
			stateFilePath: testStateFilePath,
			mockStateReader: func(t *testing.T) stateReader {
				return func(buf *bytes.Buffer) (RunnerManagers, error) {
					assert.Equal(t, testStateOutput, buf.String())

					return testRMs, nil
				}
			},
			assertRMs: func(t *testing.T, rms RunnerManagers) {
				assert.Equal(t, testRMs, rms)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			c := New(logger.New())

			require.NotNil(t, tt.mockStateReader, "mockStateReader must be defined in test definition")
			c.stateReader = tt.mockStateReader(t)

			rms, err := c.ReadStateFile(ctx, tt.stateFilePath)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, tt.assertRMs, "assertRMs must be defined in test definition")
			tt.assertRMs(t, rms)
		})
	}
}

func TestReadState(t *testing.T) {
	tests := map[string]struct {
		state       string
		assertError func(t *testing.T, err error)
		assertRMs   func(t *testing.T, rms RunnerManagers)
	}{
		"invalid state": {
			state: "invalid",
			assertError: func(t *testing.T, err error) {
				assert.Error(t, err, errParsingTerraformState)
				assert.Error(t, err, errReadingTerraformState)
			},
		},
		"invalid runner manager definition": {
			state: `{"format_version": "1", "values": {"outputs": {"grit_runner_managers": {"value": {"runner-manager-1": {"instance_name": "instance-name"}}}}}}`,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errReadingRunnerManager)
				assert.ErrorIs(t, err, errMissingRequiredEntry)
			},
		},
		"valid runner manager definition": {
			state: `{"format_version": "1", "values": {"outputs": {"grit_runner_managers": {"value": {"runner-manager-1": {"instance_name": "instance-name", "address": "address", "wrapper_address": "wrapper-address"}}}}}}`,
			assertRMs: func(t *testing.T, rms RunnerManagers) {
				if assert.Contains(t, rms, "runner-manager-1") {
					rm := rms["runner-manager-1"]
					assert.Equal(t, "instance-name", rm.InstanceName)
					assert.Equal(t, "address", rm.Address)
					assert.Equal(t, "wrapper-address", rm.WrapperAddress)
				}
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			rms, err := readState(bytes.NewBufferString(tt.state))

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, tt.assertRMs, "assertRMs must not be defined in test definition")
			tt.assertRMs(t, rms)
		})
	}
}
