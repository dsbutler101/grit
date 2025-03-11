package ssh

import (
	"context"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

type ProxyCommandDialer struct {
	sc        *streamCommand
	sshClient *ssh.Client

	cfg       Config
	clientCfg *ssh.ClientConfig
}

func NewProxyCommandDialer(cfg Config, def CommandDef) (*ProxyCommandDialer, error) {
	clientCfg, err := cfg.SSHClientConfig()
	if err != nil {
		return nil, err
	}

	sc := newStreamCommand(def.Cmd, def.Args...)

	pcd := &ProxyCommandDialer{
		sc:        sc,
		cfg:       cfg,
		clientCfg: clientCfg,
	}

	return pcd, nil
}

func (d *ProxyCommandDialer) Start(ctx context.Context) error {
	err := d.sc.Start(ctx)
	if err != nil {
		return err
	}

	conn, channels, requests, err := ssh.NewClientConn(d.sc.Connection(), d.cfg.Address, d.clientCfg)
	if err != nil {
		return fmt.Errorf("starting tunneled SSH connection: %w", err)
	}

	d.sshClient = ssh.NewClient(conn, channels, requests)

	return nil
}

func (d *ProxyCommandDialer) Wait() error {
	return d.sc.Wait()
}

func (d *ProxyCommandDialer) Dial(network string, address string) (net.Conn, error) {
	if d.sshClient == nil {
		return nil, fmt.Errorf("ProxyCommand not initialized")
	}

	return d.sshClient.Dial(network, address)
}

func (d *ProxyCommandDialer) Close() error {
	if d.sshClient != nil {
		err := d.sshClient.Close()
		if err != nil {
			return fmt.Errorf("closing inner connection: %w", err)
		}
	}

	err := d.sc.Kill()
	if err != nil {
		return fmt.Errorf("terminating ProxyCommand: %w", err)
	}

	return nil
}
