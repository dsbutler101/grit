package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDialerFactory_Create(t *testing.T) {
	testDirectDialer := new(DirectDialer)
	testProxyJumpDialer := new(ProxyJumpDialer)
	testProxyCommandDialer := new(ProxyCommandDialer)
	testCommandDialer := new(CommandDialer)

	testAddress := "test-address"
	testAddress2 := "test-address-2"
	testUsername := "test-username"
	testUsername2 := "test-username-2"
	testPrivateKeyPem := []byte("test-private-key-pem")
	testPrivateKeyPem2 := []byte("test-private-key-pem-2")
	testCommand := "test-command"

	testCommandDef := CommandDef{
		Cmd:  "test-cmd",
		Args: []string{"test-arg"},
	}

	testKeyFile := filepath.Join(t.TempDir(), "test-key.pem")
	require.NoError(t, os.WriteFile(testKeyFile, testPrivateKeyPem2, 0o600))

	defaultTestTarget := TargetDef{
		Host: TargetHostDef{
			Address:       testAddress,
			Username:      testUsername,
			PrivateKeyPem: testPrivateKeyPem,
		},
	}

	defaultNewCommandDefMock := func(t *testing.T, cmd string, target TargetDef) (CommandDef, error) {
		assert.Equal(t, testCommand, cmd)
		assert.Equal(t, defaultTestTarget, target)

		return testCommandDef, nil
	}

	failingNewCommandDefMock := func(t *testing.T, cmd string, target TargetDef) (CommandDef, error) {
		assert.Equal(t, testCommand, cmd)
		assert.Equal(t, defaultTestTarget, target)

		return CommandDef{}, assert.AnError
	}

	tests := map[string]struct {
		flags                           Flags
		target                          TargetDef
		mockNewCommandDef               func(t *testing.T, cmd string, target TargetDef) (CommandDef, error)
		assertError                     func(t *testing.T, err error)
		expectedDialer                  Dialer
		assertDirectDialerCreator       func(t *testing.T, cfg Config)
		assertProxyJumpDialerCreator    func(t *testing.T, targetCfg Config, proxyCfg Config)
		assertProxyCommandDialerCreator func(t *testing.T, cfg Config, def CommandDef)
		assertCommandDialerCreator      func(t *testing.T, def CommandDef)
	}{
		"direct dialer creation": {
			target:         defaultTestTarget,
			expectedDialer: testDirectDialer,
			assertDirectDialerCreator: func(t *testing.T, cfg Config) {
				assert.Equal(t, testAddress, cfg.Address)
				assert.Equal(t, testUsername, cfg.Username)
				assert.Equal(t, testPrivateKeyPem, cfg.KeyPemBytes)
			},
		},
		"proxy jump dialer missing username": {
			flags: Flags{
				ProxyJump: ProxyJumpFlags{
					Address: testAddress2,
				},
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errProxyJumpMissingUsername)
			},
		},
		"proxy jump dialer missing key file": {
			flags: Flags{
				ProxyJump: ProxyJumpFlags{
					Address:  testAddress2,
					Username: testUsername2,
				},
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errProxyJumpMissingKeyFile)
			},
		},
		"proxy jump dialer key file reading error": {
			flags: Flags{
				ProxyJump: ProxyJumpFlags{
					Address:  testAddress2,
					Username: testUsername2,
					KeyFile:  "unknown-file",
				},
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errProxyJumpReadingKeyFile)
			},
		},
		"proxy jump dialer creation": {
			flags: Flags{
				ProxyJump: ProxyJumpFlags{
					Address:  testAddress2,
					Username: testUsername2,
					KeyFile:  testKeyFile,
				},
			},
			target:         defaultTestTarget,
			expectedDialer: testProxyJumpDialer,
			assertProxyJumpDialerCreator: func(t *testing.T, targetCfg Config, proxyCfg Config) {
				assert.Equal(t, testAddress, targetCfg.Address)
				assert.Equal(t, testUsername, targetCfg.Username)
				assert.Equal(t, testPrivateKeyPem, targetCfg.KeyPemBytes)

				assert.Equal(t, testAddress2, proxyCfg.Address)
				assert.Equal(t, testUsername2, proxyCfg.Username)
				assert.Equal(t, testPrivateKeyPem2, proxyCfg.KeyPemBytes)
			},
		},
		"proxy command dialer newCommandDef failure": {
			flags: Flags{
				ProxyCommand: testCommand,
			},
			target:            defaultTestTarget,
			mockNewCommandDef: failingNewCommandDefMock,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, errCommandDefCreation)
			},
		},
		"proxy command dialer creation": {
			flags: Flags{
				ProxyCommand: testCommand,
			},
			target:            defaultTestTarget,
			mockNewCommandDef: defaultNewCommandDefMock,
			expectedDialer:    testProxyCommandDialer,
			assertProxyCommandDialerCreator: func(t *testing.T, cfg Config, def CommandDef) {
				assert.Equal(t, testAddress, cfg.Address)
				assert.Equal(t, testUsername, cfg.Username)
				assert.Equal(t, testPrivateKeyPem, cfg.KeyPemBytes)

				assert.Equal(t, testCommandDef, def)
			},
		},
		"command dialer newCommandDef failure": {
			flags: Flags{
				Command: testCommand,
			},
			mockNewCommandDef: failingNewCommandDefMock,
			target:            defaultTestTarget,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, errCommandDefCreation)
			},
		},
		"command dialer creation": {
			flags: Flags{
				Command: testCommand,
			},
			target:            defaultTestTarget,
			mockNewCommandDef: defaultNewCommandDefMock,
			expectedDialer:    testCommandDialer,
			assertCommandDialerCreator: func(t *testing.T, def CommandDef) {
				assert.Equal(t, testCommandDef, def)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			df := NewDialerFactory()

			df.createDirectDialer = func(cfg Config) (*DirectDialer, error) {
				require.NotNil(t, tt.assertDirectDialerCreator, "assertDirectDialerCreator must be defined in test definition")
				tt.assertDirectDialerCreator(t, cfg)

				return testDirectDialer, nil
			}

			df.createProxyJumpDialer = func(targetCfg Config, proxyCfg Config) (*ProxyJumpDialer, error) {
				require.NotNil(t, tt.assertProxyJumpDialerCreator, "assertProxyJumpDialerCreator must be defined in test definition")
				tt.assertProxyJumpDialerCreator(t, targetCfg, proxyCfg)

				return testProxyJumpDialer, nil
			}

			df.createProxyCommandDialer = func(cfg Config, def CommandDef) (*ProxyCommandDialer, error) {
				require.NotNil(t, tt.assertProxyCommandDialerCreator, "assertProxyCommandDialerCreator must be defined in test definition")
				tt.assertProxyCommandDialerCreator(t, cfg, def)

				return testProxyCommandDialer, nil
			}

			df.createCommandDialer = func(def CommandDef) *CommandDialer {
				require.NotNil(t, tt.assertCommandDialerCreator, "assertCommandDialerCreator must be defined in test definition")
				tt.assertCommandDialerCreator(t, def)

				return testCommandDialer
			}

			df.newCommandDef = func(cmd string, target TargetDef) (CommandDef, error) {
				require.NotNil(t, tt.mockNewCommandDef, "mockNewCommandDef must be defined in test definition")
				return tt.mockNewCommandDef(t, cmd, target)
			}

			d, err := df.Create(tt.flags, tt.target)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedDialer, d)
		})
	}
}

func TestNewCommandDef(t *testing.T) {
	testAddress := "test-address"
	testUsername := "test-username"
	testPrivateKeyPem := []byte("test-private-key")

	testGRPCNetwork := "test-grpc-network"
	testGRPCAddress := "test-grpc-address"

	testTargetDef := TargetDef{
		Host: TargetHostDef{
			Address:       testAddress,
			Username:      testUsername,
			PrivateKeyPem: testPrivateKeyPem,
		},
		GRPCServer: TargetGRPCServerDef{
			Network: testGRPCNetwork,
			Address: testGRPCAddress,
		},
	}

	tests := map[string]struct {
		command     string
		target      TargetDef
		assertError func(t *testing.T, err error)
		assertDef   func(t *testing.T, d CommandDef)
	}{
		"empty command": {
			command: "",
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errCommandDefEmptyCommand)
			},
		},
		"template parsing error": {
			command: "{{",
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errCommandDefParsingCommandTemplate)
			},
		},
		"template execution error": {
			command: "{{ .UnknownField }}",
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errCommandDefExecutingCommandTemplate)
			},
		},
		"command definition empty after evaluation": {
			command: "{{ .Host.Address }}",
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errCommandDefInvalidCommandDefinition)
			},
		},
		"command without arguments": {
			command: "cmd",
			target:  testTargetDef,
			assertDef: func(t *testing.T, d CommandDef) {
				assert.Equal(t, "cmd", d.Cmd)
				assert.Empty(t, d.Args)
			},
		},
		"command with arguments": {
			command: "cmd {{ .Host.Address }} {{ .Host.Username }} {{ .GRPCServer.Network }}:{{ .GRPCServer.Address}}",
			target:  testTargetDef,
			assertDef: func(t *testing.T, d CommandDef) {
				assert.Equal(t, "cmd", d.Cmd)
				if assert.Len(t, d.Args, 3) {
					assert.Equal(t, d.Args[0], testAddress)
					assert.Equal(t, d.Args[1], testUsername)
					assert.Equal(t, d.Args[2], fmt.Sprintf("%s:%s", testGRPCNetwork, testGRPCAddress))
				}
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			def, err := newCommandDef(tt.command, tt.target)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, tt.assertDef, "assertDef must be defined in test definition")
			tt.assertDef(t, def)
		})
	}
}
