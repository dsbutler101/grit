package ssh

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api/client"
	"golang.org/x/crypto/ssh"
)

const (
	testSSHUsername          = "ssh-test-user"
	testSSHUserPrivateKeyPem = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAaAAAABNlY2RzYS
1zaGEyLW5pc3RwMjU2AAAACG5pc3RwMjU2AAAAQQTwQiqRdyXcKM7HXGYNlBZiAqxeRSOL
n2+f6QXC4Vdqf9DsHHhvGVfBAOJtR0tymJCmy/NkRcUUSM/XJX6KLnd7AAAAsJGc8rGRnP
KxAAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBPBCKpF3Jdwozsdc
Zg2UFmICrF5FI4ufb5/pBcLhV2p/0OwceG8ZV8EA4m1HS3KYkKbL82RFxRRIz9clfooud3
sAAAAgZTtaOL0rmx/wHQYSCAWKAcvIEunoC/RLaZU+3Ocm2FcAAAARdG1hY3p1a2luQHBl
Z2FzdXMBAgMEBQYH
-----END OPENSSH PRIVATE KEY-----`
	testSSHHostPrivateKeyPem = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAaAAAABNlY2RzYS
1zaGEyLW5pc3RwMjU2AAAACG5pc3RwMjU2AAAAQQRle4k1J9NX3KZAL9rroLLBqpJ6MOrg
vcwsXZ+gEmYBUP5uAGBTLD6nWuHhR6IGCkpIGb9+kmq7AhTZl4uGXgbZAAAAsFtD+kNbQ/
pDAAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBGV7iTUn01fcpkAv
2uugssGqknow6uC9zCxdn6ASZgFQ/m4AYFMsPqda4eFHogYKSkgZv36SarsCFNmXi4ZeBt
kAAAAgaST4xruswXVEl+ecesgj97/ZIrktMtAfdvZDpoQoStUAAAARdG1hY3p1a2luQHBl
Z2FzdXMBAgMEBQYH
-----END OPENSSH PRIVATE KEY-----`
)

type closer func()

func startTestSSHServer(t *testing.T) (Config, closer) {
	testSSHUserPrivateKeyPemBytes := []byte(testSSHUserPrivateKeyPem)

	testUserKey, err := ssh.ParsePrivateKey(testSSHUserPrivateKeyPemBytes)
	require.NoError(t, err)

	testHostKey, err := ssh.ParsePrivateKey([]byte(testSSHHostPrivateKeyPem))
	require.NoError(t, err)

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(m ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if string(key.Marshal()) != string(testUserKey.PublicKey().Marshal()) {
				return nil, fmt.Errorf("unknown public key for %q", m.User())
			}

			if m.User() != testSSHUsername {
				return nil, fmt.Errorf("unknown user %q", m.User())
			}

			perm := &ssh.Permissions{
				Extensions: map[string]string{
					"pubkey-fp": ssh.FingerprintSHA256(key),
				},
			}

			return perm, nil
		},
	}
	config.AddHostKey(testHostKey)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	closeFn := func() {
		assert.NoError(t, listener.Close())
	}

	go handleTestSSHServerConnection(t, listener, config)

	cfg := Config{
		Address:     listener.Addr().String(),
		Username:    testSSHUsername,
		KeyPemBytes: testSSHUserPrivateKeyPemBytes,
	}

	return cfg, closeFn
}

func handleTestSSHServerConnection(t *testing.T, listener net.Listener, config *ssh.ServerConfig) {
	tcpConn, err := listener.Accept()
	require.NoError(t, err)
	defer func() {
		_ = tcpConn.Close()
	}()

	conn, channels, requests, err := ssh.NewServerConn(tcpConn, config)
	require.NoError(t, err)

	t.Logf("logged in as %s with key %s", conn.User(), conn.Permissions.Extensions["pubkey-fp"])

	wg := new(sync.WaitGroup)
	defer wg.Wait()

	wg.Add(1)
	go func() {
		ssh.DiscardRequests(requests)
		wg.Done()
	}()

	for newChannel := range channels {
		testSSHServerHandleNewChannel(t, wg, newChannel)
	}
}

func testSSHServerHandleNewChannel(t *testing.T, wg *sync.WaitGroup, newChannel ssh.NewChannel) {
	handlers := map[string]func(t *testing.T, channel ssh.Channel, extraData []byte){
		"direct-streamlocal@openssh.com": testSSHServerUnixForwardChannel,
		"direct-tcpip":                   testSSHServerTCPIPForwardChannel,
	}

	handler, ok := handlers[newChannel.ChannelType()]
	if !ok {
		assert.NoError(t, newChannel.Reject(ssh.Prohibited, "connection type not supported"))
		assert.Fail(t, "unexpected channel type: %s", newChannel.ChannelType())
		return
	}

	channel, requests, err := newChannel.Accept()
	require.NoError(t, err)
	defer channel.Close()

	wg.Add(1)
	go func() {
		ssh.DiscardRequests(requests)
		wg.Done()
	}()

	handler(t, channel, newChannel.ExtraData())
}

func testSSHServerUnixForwardChannel(t *testing.T, channel ssh.Channel, extraData []byte) {
	type socketForwardMsg struct {
		SocketPath     string
		ReservedStr    string
		ReservedUint32 uint32
	}

	var msg socketForwardMsg
	require.NoError(t, ssh.Unmarshal(extraData, &msg))

	testSSHServerProxyTo(t, channel, "unix", msg.SocketPath)
}

func testSSHServerProxyTo(t *testing.T, channel ssh.Channel, network string, address string) {
	proxyConn, err := net.Dial(network, address)
	require.NoError(t, err)
	defer proxyConn.Close()

	t.Log("opened proxy to", address)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		_, err := io.Copy(proxyConn, channel)
		assert.NoError(t, err)
	}()

	go func() {
		defer wg.Done()

		_, err := io.Copy(channel, proxyConn)
		assert.NoError(t, err)
	}()

	wg.Wait()
}

func testSSHServerTCPIPForwardChannel(t *testing.T, channel ssh.Channel, extraData []byte) {
	type socketForwardMsg struct {
		HostToConnect       string
		PortToConnect       uint32
		OriginatorIPAddress string
		OriginatorPort      uint32
	}

	var msg socketForwardMsg
	require.NoError(t, ssh.Unmarshal(extraData, &msg))

	testSSHServerProxyTo(t, channel, "tcp", fmt.Sprintf("%s:%d", msg.HostToConnect, msg.PortToConnect))
}

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

func withCompiledNC(t *testing.T, fn func(t *testing.T, ncPath string)) {
	sourcePath := "." + string(filepath.Separator) + filepath.Join("testdata", "nc")
	binPath := filepath.Join(t.TempDir(), "nc.exe")

	cmd := exec.Command("go", "build", "-o", binPath, sourcePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	require.NoError(t, err, "Compiling nc binary")

	fn(t, binPath)
}

func withSSHServerCfg(t *testing.T, fn func(t *testing.T, cfg Config)) {
	cfg, closeFn := startTestSSHServer(t)
	defer closeFn()

	t.Logf("Started SSH server at %s", cfg.Address)

	fn(t, cfg)
}

func runTestCall(t *testing.T, dialer client.Dialer, socketPath string) {
	httpCli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return dialer("unix", socketPath)
			},
		},
	}

	resp, err := httpCli.Get("http://ignore/test")
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	t.Log("client received response", resp.Status)

	_, err = io.Copy(io.Discard, resp.Body)
	assert.NoError(t, err)
}

func withSocket(t *testing.T, fn func(t *testing.T, socketPath string)) {
	dir := t.TempDir()
	socketPath := filepath.Join(dir, "test.sock")

	listener, err := net.Listen("unix", socketPath)
	require.NoError(t, err)

	t.Log("created socket at", socketPath)

	server := &httptest.Server{
		Listener: listener,
		Config: &http.Server{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Logf("server received request to %s from %s", r.RequestURI, r.RemoteAddr)
				if r.URL.Path == "/test" {
					w.WriteHeader(http.StatusOK)
					return
				}

				w.WriteHeader(http.StatusNotFound)
			}),
		},
	}
	server.Start()
	defer server.Close()

	t.Log("started server on the socket")

	fn(t, socketPath)
}

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
