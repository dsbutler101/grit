package wrapper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/logger"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api/client"
)

func TestMux_Execute(t *testing.T) {
	testTarget := "/test/target/dir"
	testCallbackFn := func(ctx context.Context, c *Client) error {
		return nil
	}

	testRMAlias1 := "runner-manager-1"
	testRM1 := terraform.RunnerManager{
		InstanceName: "test-instance",
	}
	testRunnerManagers := terraform.RunnerManagers{
		testRMAlias1: testRM1,
	}

	tests := map[string]struct {
		tfFlags                     terraform.Flags
		prepareTFClientMock         func(ctx context.Context, client *mockTfClient)
		prepareRMHandlerFactoryMock func(t *testing.T, ctx context.Context, factory *mockRmHandlerFactory)
		callbackFn                  callbackFn
		cancelContext               bool
		assertError                 func(t *testing.T, err error)
	}{
		"nil callback function": {
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrNilCallbackFunction)
			},
		},
		"runner managers listing failure": {
			tfFlags: terraform.Flags{
				Target: testTarget,
			},
			prepareTFClientMock: func(ctx context.Context, client *mockTfClient) {
				client.EXPECT().ReadStateDir(ctx, testTarget).Return(nil, assert.AnError)
			},
			callbackFn: testCallbackFn,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		"no runner managers on the list": {
			tfFlags: terraform.Flags{
				TargetStateFile: testTarget,
			},
			prepareTFClientMock: func(ctx context.Context, client *mockTfClient) {
				client.EXPECT().ReadStateFile(ctx, testTarget).Return(nil, nil)
			},
			callbackFn: testCallbackFn,
		},
		"no runner managers on the list and context cancelled": {
			tfFlags: terraform.Flags{
				TargetStateFile: testTarget,
			},
			prepareTFClientMock: func(ctx context.Context, client *mockTfClient) {
				client.EXPECT().ReadStateFile(ctx, testTarget).Return(nil, nil)
			},
			callbackFn:    testCallbackFn,
			cancelContext: true,
		},
		"runner manager handled successfully": {
			tfFlags: terraform.Flags{
				TargetStateFile: testTarget,
			},
			prepareTFClientMock: func(ctx context.Context, client *mockTfClient) {
				client.EXPECT().
					ReadStateFile(ctx, testTarget).
					Return(testRunnerManagers, nil)
			},
			prepareRMHandlerFactoryMock: func(t *testing.T, ctx context.Context, factory *mockRmHandlerFactory) {
				handler := newMockRmHandler(t)
				factory.EXPECT().new(testRMAlias1).Return(handler)

				handler.EXPECT().handle(ctx, testRM1, mock.AnythingOfType("callbackFn")).Return(nil)
			},
			callbackFn: testCallbackFn,
		},
		"runner manager failed": {
			tfFlags: terraform.Flags{
				TargetStateFile: testTarget,
			},
			prepareTFClientMock: func(ctx context.Context, client *mockTfClient) {
				client.EXPECT().
					ReadStateFile(ctx, testTarget).
					Return(testRunnerManagers, nil)
			},
			prepareRMHandlerFactoryMock: func(t *testing.T, ctx context.Context, factory *mockRmHandlerFactory) {
				handler := newMockRmHandler(t)
				factory.EXPECT().new(testRMAlias1).Return(handler)

				handler.EXPECT().handle(ctx, testRM1, mock.AnythingOfType("callbackFn")).Return(assert.AnError)
			},
			callbackFn: testCallbackFn,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			sshFlags := ssh.Flags{
				Username: "",
				KeyFile:  "",
				ProxyJump: ssh.ProxyJumpFlags{
					Address:  "",
					Username: "",
					KeyFile:  "",
				},
				ProxyCommand: "",
				Command:      "",
			}

			wrapperFlags := Flags{
				ConnectionTimeout: DefaultTimeout,
			}

			tfClientMock := newMockTfClient(t)
			if tt.prepareTFClientMock != nil {
				tt.prepareTFClientMock(ctx, tfClientMock)
			}

			if tt.cancelContext {
				cancelFn()
			}

			mux := NewMux(logger.New(), tfClientMock, tt.tfFlags, sshFlags, wrapperFlags)

			rmHandlerFactoryMock := newMockRmHandlerFactory(t)
			mux.rmHandlerFactory = rmHandlerFactoryMock

			if tt.prepareRMHandlerFactoryMock != nil {
				tt.prepareRMHandlerFactoryMock(t, ctx, rmHandlerFactoryMock)
			}

			err := mux.Execute(ctx, tt.callbackFn)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestGetWrapperAddress(t *testing.T) {
	tests := map[string]struct {
		address       string
		assertAddress func(t *testing.T, address wrapperAddress)
		assertError   func(t *testing.T, err error)
	}{
		"invalid address format": {
			address: "{some-address:1234",
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrWrapperAddressParsing)
			},
		},
		"invalid address network": {
			address: "some-address:1234",
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrInvalidWrapperAddressNetwork)
			},
		},
		"unix address": {
			address: "unix:///var/run/wrapper.sock",
			assertAddress: func(t *testing.T, address wrapperAddress) {
				assert.Equal(t, "unix", address.network)
				assert.Equal(t, "/var/run/wrapper.sock", address.address)
			},
		},
		"tcp address": {
			address: "tcp://127.0.0.1:1234",
			assertAddress: func(t *testing.T, address wrapperAddress) {
				assert.Equal(t, "tcp", address.network)
				assert.Equal(t, "127.0.0.1:1234", address.address)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			address, err := getWrapperAddress(terraform.RunnerManager{
				WrapperAddress: tt.address,
			})

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, tt.assertAddress, "assertAddress must be defined in the test definition")

			tt.assertAddress(t, address)
		})
	}
}

func TestGetUsername(t *testing.T) {
	testUsername := "test-username"
	testUsername2 := "test-username-2"

	tests := map[string]struct {
		username         string
		sshFlags         ssh.Flags
		assertError      func(t *testing.T, err error)
		expectedUsername string
	}{
		"username provided with runner manager details": {
			username:         testUsername,
			expectedUsername: testUsername,
		},
		"username provided with ssh flags": {
			sshFlags: ssh.Flags{
				Username: testUsername,
			},
			expectedUsername: testUsername,
		},
		"username provided with both runner manager details and ssh flags": {
			username: testUsername,
			sshFlags: ssh.Flags{
				Username: testUsername2,
			},
			expectedUsername: testUsername2,
		},
		"username not provided at all": {
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrUsernameNotProvided)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			username, err := getUsername(terraform.RunnerManager{Username: tt.username}, tt.sshFlags)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedUsername, username)
		})
	}
}

func TestGetSSHKeyPemBytes(t *testing.T) {
	testKeyPem := "key-pem"
	testKeyPem2 := "key-pem-2"

	dir := t.TempDir()

	testKeyPem2File := filepath.Join(dir, "test.key")
	require.NoError(t, os.WriteFile(testKeyPem2File, []byte(testKeyPem2), 0600))

	tests := map[string]struct {
		sshKeyPem           string
		sshFlags            ssh.Flags
		assertError         func(t *testing.T, err error)
		expectedKeyPemBytes []byte
	}{
		"key provided with runner manager details": {
			sshKeyPem:           testKeyPem,
			expectedKeyPemBytes: []byte(testKeyPem),
		},
		"key provided with ssh flags": {
			sshFlags: ssh.Flags{
				KeyFile: testKeyPem2File,
			},
			expectedKeyPemBytes: []byte(testKeyPem2),
		},
		"key provided with both runner manager details and ssh flags": {
			sshKeyPem: testKeyPem,
			sshFlags: ssh.Flags{
				KeyFile: testKeyPem2File,
			},
			expectedKeyPemBytes: []byte(testKeyPem2),
		},
		"ssh key file from ssh flags reading failure": {
			sshFlags: ssh.Flags{
				KeyFile: filepath.Join(dir, "unknown.key"),
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrReadingSSHKeyFile)
			},
		},
		"key not provided at all": {
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrSSHKeyPemNotProvided)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			keyPemBytes, err := getSSHKeyPemBytes(terraform.RunnerManager{SSHKeyPem: tt.sshKeyPem}, tt.sshFlags)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedKeyPemBytes, keyPemBytes)
		})
	}
}

func TestMuxRunnerManagerHandler_handle(t *testing.T) {
	testProxyCommand := "proxy-command"
	testSSHFlags := ssh.Flags{
		ProxyCommand: testProxyCommand,
	}

	testAlias := "test-runner-manager-1"
	testAddress := "localhost:1234"
	testSocketNetwork := "unix"
	testSocketPath := "/var/run/wrapper.sock"
	testSocket := fmt.Sprintf("%s:%s", testSocketNetwork, testSocketPath)
	testUsername := "test-username-1"
	testKeyPem := "test-key-pem"

	testRM := terraform.RunnerManager{
		Address:        testAddress,
		WrapperAddress: testSocket,
		Username:       testUsername,
		SSHKeyPem:      testKeyPem,
	}

	noopDialerFactory := func(_ *testing.T, _ context.Context, _ *ssh.MockDialer) dialerFactory {
		return func(_ ssh.Flags, _ ssh.TargetDef) (ssh.Dialer, error) {
			return nil, nil
		}
	}

	defaultDialerFactory := func(t *testing.T, expectedCtx context.Context, dialer *ssh.MockDialer) dialerFactory {
		return func(flags ssh.Flags, target ssh.TargetDef) (ssh.Dialer, error) {
			assert.Equal(t, testSSHFlags, flags)
			assert.Equal(t, testRM.Address, target.Host.Address)
			assert.Equal(t, testRM.Username, target.Host.Username)
			assert.Equal(t, []byte(testRM.SSHKeyPem), target.Host.PrivateKeyPem)
			assert.Equal(t, testSocketNetwork, target.GRPCServer.Network)
			assert.Equal(t, testSocketPath, target.GRPCServer.Address)

			dialerClosedCh := make(chan struct{})

			dialer.EXPECT().Start(expectedCtx).Return(nil)
			dialer.EXPECT().
				Close().
				Run(func() {
					close(dialerClosedCh)
				}).
				Return(nil)
			dialer.EXPECT().
				Wait().
				Run(func() {
					select {
					case <-expectedCtx.Done():
					case <-dialerClosedCh:
					}
				}).
				Return(errors.New("logged only error"))

			return dialer, nil
		}
	}

	noopClientFactory := func(_ *testing.T, _ context.Context, _ *ssh.MockDialer, _ *Client) clientFactory {
		return func(ctx context.Context, logger *slog.Logger, dialer client.Dialer, connectionTimeout time.Duration, address string) (*Client, error) {
			return nil, nil
		}
	}

	defaultClientFactory := func(t *testing.T, expectedCtx context.Context, expectedDialer *ssh.MockDialer, cl *Client) clientFactory {
		return func(ctx context.Context, logger *slog.Logger, dialer client.Dialer, connectionTimeout time.Duration, address string) (*Client, error) {
			assert.Equal(t, expectedCtx, ctx)

			return cl, nil
		}
	}

	noopCallbackFn := func(_ *testing.T, _ context.Context, _ *Client) callbackFn {
		return func(_ context.Context, _ *Client) error {
			return nil
		}
	}

	tests := map[string]struct {
		sshFlags      ssh.Flags
		rm            terraform.RunnerManager
		dialerFactory func(t *testing.T, expectedCtx context.Context, dialerMock *ssh.MockDialer) dialerFactory
		clientFactory func(t *testing.T, expectedCtx context.Context, expectedDialer *ssh.MockDialer, client *Client) clientFactory
		callbackFn    func(t *testing.T, expectedCtx context.Context, expectedClient *Client) callbackFn
		assertError   func(t *testing.T, err error)
	}{
		"runner manager address failure": {
			rm: terraform.RunnerManager{
				WrapperAddress: "unknown://localhost:1234",
			},
			dialerFactory: noopDialerFactory,
			clientFactory: noopClientFactory,
			callbackFn:    noopCallbackFn,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrInvalidWrapperAddressNetwork)
			},
		},
		"username failure": {
			rm: terraform.RunnerManager{
				WrapperAddress: testSocket,
			},
			dialerFactory: noopDialerFactory,
			clientFactory: noopClientFactory,
			callbackFn:    noopCallbackFn,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrUsernameNotProvided)
			},
		},
		"ssh key failure": {
			rm: terraform.RunnerManager{
				WrapperAddress: testSocket,
				Username:       testUsername,
			},
			dialerFactory: noopDialerFactory,
			clientFactory: noopClientFactory,
			callbackFn:    noopCallbackFn,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrSSHKeyPemNotProvided)
			},
		},
		"dialer factory error": {
			rm: testRM,
			dialerFactory: func(_ *testing.T, _ context.Context, _ *ssh.MockDialer) dialerFactory {
				return func(_ ssh.Flags, _ ssh.TargetDef) (ssh.Dialer, error) {
					return nil, assert.AnError
				}
			},
			clientFactory: noopClientFactory,
			callbackFn:    noopCallbackFn,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrSSHDialerFabrication)
			},
		},
		"dialer start failure": {
			rm: testRM,
			dialerFactory: func(t *testing.T, expectedCtx context.Context, dialer *ssh.MockDialer) dialerFactory {
				return func(_ ssh.Flags, _ ssh.TargetDef) (ssh.Dialer, error) {
					dialer.EXPECT().Start(expectedCtx).Return(assert.AnError)
					dialer.EXPECT().Close().Return(errors.New("logged only error"))

					return dialer, nil
				}
			},
			clientFactory: noopClientFactory,
			callbackFn:    noopCallbackFn,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrSSHDialerStart)
			},
		},
		"client creation failure": {
			sshFlags:      testSSHFlags,
			rm:            testRM,
			dialerFactory: defaultDialerFactory,
			clientFactory: func(t *testing.T, expectedCtx context.Context, expectedDialer *ssh.MockDialer, _ *Client) clientFactory {
				return func(ctx context.Context, logger *slog.Logger, dialer client.Dialer, connectionTimeout time.Duration, address string) (*Client, error) {
					assert.Equal(t, expectedCtx, ctx)

					return nil, assert.AnError
				}
			},
			callbackFn: noopCallbackFn,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrWrapperConnect)
			},
		},
		"callback function failure": {
			sshFlags:      testSSHFlags,
			rm:            testRM,
			dialerFactory: defaultDialerFactory,
			clientFactory: defaultClientFactory,
			callbackFn: func(t *testing.T, expectedCtx context.Context, expectedClient *Client) callbackFn {
				return func(ctx context.Context, client *Client) error {
					assert.Equal(t, expectedCtx, ctx)
					assert.Equal(t, expectedClient, client)

					return assert.AnError
				}
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrWrapperExecution)
			},
		},
		"successful execution": {
			sshFlags:      testSSHFlags,
			rm:            testRM,
			dialerFactory: defaultDialerFactory,
			clientFactory: defaultClientFactory,
			callbackFn: func(t *testing.T, expectedCtx context.Context, expectedClient *Client) callbackFn {
				return func(ctx context.Context, client *Client) error {
					assert.Equal(t, expectedCtx, ctx)
					assert.Equal(t, expectedClient, client)

					return nil
				}
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			ctx, cancelFn := context.WithCancel(context.Background())
			defer cancelFn()

			require.NotNil(t, tt.dialerFactory, "dialerFactory must be defined in test definition")
			require.NotNil(t, tt.clientFactory, "clientFactory must be defined in test definition")

			dialer := ssh.NewMockDialer(t)
			cl := new(Client)

			factory := &muxRunnerManagerHandlerFactory{
				logger:        logger.New(),
				sshFlags:      tt.sshFlags,
				wrapperFlags:  Flags{},
				dialerFactory: tt.dialerFactory(t, ctx, dialer),
				clientFactory: tt.clientFactory(t, ctx, dialer, cl),
			}

			require.NotNil(t, tt.callbackFn, "callbackFn must be defined in test definition")

			handler := factory.new(testAlias)
			err := handler.handle(ctx, tt.rm, tt.callbackFn(t, ctx, cl))

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
