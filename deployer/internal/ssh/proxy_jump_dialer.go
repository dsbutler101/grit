package ssh

import (
	"context"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

type ProxyJumpDialer struct {
	proxyCli  *ssh.Client
	targetCli *ssh.Client
}

func NewProxyJumpDialer(targetCfg Config, proxyCfg Config) (*ProxyJumpDialer, error) {
	proxyCliCfg, err := proxyCfg.SSHClientConfig()
	if err != nil {
		return nil, fmt.Errorf("preparing SSH client config for Proxy node: %w", err)
	}

	clientCfg, err := targetCfg.SSHClientConfig()
	if err != nil {
		return nil, fmt.Errorf("preparing SSH client config for Address node: %w", err)
	}

	proxyCli, err := ssh.Dial(proxyCfg.Network(), proxyCfg.Address, proxyCliCfg)
	if err != nil {
		return nil, fmt.Errorf("dialing SSH Proxy node: %w", err)
	}

	targetConn, err := proxyCli.Dial(targetCfg.Network(), targetCfg.Address)
	if err != nil {
		return nil, fmt.Errorf("dialing SSH Address node: %w", err)
	}

	targetSSHConn, channels, requests, err := ssh.NewClientConn(targetConn, targetCfg.Address, clientCfg)
	if err != nil {
		return nil, fmt.Errorf("creating SSH client connection for Address node: %w", err)
	}

	dialer := &ProxyJumpDialer{
		proxyCli:  proxyCli,
		targetCli: ssh.NewClient(targetSSHConn, channels, requests),
	}

	return dialer, nil
}

func (d *ProxyJumpDialer) Start(_ context.Context) error {
	return nil
}

func (d *ProxyJumpDialer) Wait() error {
	return nil
}

func (d *ProxyJumpDialer) Dial(network string, address string) (net.Conn, error) {
	return d.targetCli.Dial(network, address)
}

func (d *ProxyJumpDialer) Close() error {
	err := d.targetCli.Close()
	if err != nil {
		return fmt.Errorf("closing SSH connection to Address node: %w", err)
	}

	err = d.proxyCli.Close()
	if err != nil {
		return fmt.Errorf("closing SSH connection to Proxy node: %w", err)
	}

	return nil
}
