package ssh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"text/template"
)

type TargetDef struct {
	Host       TargetHostDef
	GRPCServer TargetGRPCServerDef
}

type TargetHostDef struct {
	Address       string
	Username      string
	PrivateKeyPem []byte
}

type TargetGRPCServerDef struct {
	Network string
	Address string
}

//go:generate mockery --name=Dialer --inpackage --with-expecter
type Dialer interface {
	Dial(network string, address string) (net.Conn, error)
	Start(ctx context.Context) error
	Wait() error
	Close() error
}

func NewDialer(flags Flags, target TargetDef) (Dialer, error) {
	if flags.Command != "" {
		return newCommandDialer(flags, target)
	}

	if flags.ProxyCommand != "" {
		return newProxyCommandDialer(flags, target)
	}

	if flags.ProxyJump.Address != "" {
		return newProxyJumpDialer(flags, target)
	}

	return newDirectDialer(flags, target)
}

func newCommandDialer(flags Flags, target TargetDef) (Dialer, error) {
	def, err := newCommandDef(flags.Command, target)
	if err != nil {
		return nil, err
	}

	return NewCommandDialer(def), nil
}

func newCommandDef(cmd string, target TargetDef) (CommandDef, error) {
	var d CommandDef

	if cmd == "" {
		return d, errors.New("command can't be empty")
	}

	tpl, err := template.New("command").Parse(cmd)
	if err != nil {
		return d, fmt.Errorf("parsing command: %w", err)
	}

	buf := bytes.NewBuffer(nil)
	err = tpl.Execute(buf, target)
	if err != nil {
		return d, fmt.Errorf("executing command template: %w", err)
	}

	parts := strings.Split(buf.String(), " ")
	if len(parts) < 1 {
		return d, errors.New("command did not contain a proper definition")
	}

	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	d.Cmd = parts[0]
	d.Args = args

	return d, nil
}

func newProxyCommandDialer(flags Flags, target TargetDef) (Dialer, error) {
	def, err := newCommandDef(flags.ProxyCommand, target)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Address:     target.Host.Address,
		Username:    target.Host.Username,
		KeyPemBytes: target.Host.PrivateKeyPem,
	}

	return NewProxyCommandDialer(cfg, def)
}

func newProxyJumpDialer(flags Flags, target TargetDef) (Dialer, error) {
	if flags.ProxyJump.Username == "" {
		return nil, fmt.Errorf("SSH Proxy Jump username must be specified")
	}

	if flags.ProxyJump.KeyFile == "" {
		return nil, fmt.Errorf("SSH Proxy Jump key file must be specified")
	}

	keyPem, err := os.ReadFile(flags.ProxyJump.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("reading SSH Proxy Jump key file: %w", err)
	}

	proxyCfg := Config{
		Address:     flags.ProxyJump.Address,
		Username:    flags.ProxyJump.Username,
		KeyPemBytes: keyPem,
	}

	targetCfg := Config{
		Address:     target.Host.Address,
		Username:    target.Host.Username,
		KeyPemBytes: target.Host.PrivateKeyPem,
	}

	return NewProxyJumpDialer(targetCfg, proxyCfg)
}

func newDirectDialer(flags Flags, target TargetDef) (Dialer, error) {
	targetCfg := Config{
		Address:     target.Host.Address,
		Username:    target.Host.Username,
		KeyPemBytes: target.Host.PrivateKeyPem,
	}

	return NewDirectDialer(targetCfg)
}
