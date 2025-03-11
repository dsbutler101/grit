package ssh

import (
	"context"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

type DirectDialer struct {
	cli *ssh.Client
}

func NewDirectDialer(cfg Config) (*DirectDialer, error) {
	clientCfg, err := cfg.SSHClientConfig()
	if err != nil {
		return nil, fmt.Errorf("preparing SSH client config: %w", err)
	}

	cli, err := ssh.Dial(cfg.Network(), cfg.Address, clientCfg)
	if err != nil {
		return nil, fmt.Errorf("dialing SSH: %w", err)
	}

	dialer := &DirectDialer{
		cli: cli,
	}

	return dialer, nil
}

func (d *DirectDialer) Start(_ context.Context) error {
	return nil
}

func (d *DirectDialer) Wait() error {
	return nil
}

func (d *DirectDialer) Dial(network string, address string) (net.Conn, error) {
	return d.cli.Dial(network, address)
}

func (d *DirectDialer) Close() error {
	return d.cli.Close()
}
