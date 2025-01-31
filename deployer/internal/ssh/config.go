package ssh

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	sshNetwork = "tcp"

	sshClientConnectionTimeout = 30 * time.Second
)

type Config struct {
	// Address holds the address of the SSH node. For simplicity, we assume
	// that SSH is always going to be reached through a TCP port, so the
	// Address field expects a `hostname:port_number` or `ip_address:port_number`
	// format.
	Address     string
	Username    string
	KeyPemBytes []byte
}

func (c Config) SSHClientConfig() (*ssh.ClientConfig, error) {
	key, err := ssh.ParsePrivateKey(c.KeyPemBytes)
	if err != nil {
		return nil, fmt.Errorf("parsing private key bytes: %w", err)
	}

	clientCfg := &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(key)},

		// This may be the easiest way to move forward with the initial implementation,
		// but at some moment we should make it configurable. User should be able to
		// choose between ignoring host keys or configuring expected host keys that
		// should be enforced.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		Timeout: sshClientConnectionTimeout,
	}

	return clientCfg, nil
}

func (c Config) Network() string {
	return sshNetwork
}
