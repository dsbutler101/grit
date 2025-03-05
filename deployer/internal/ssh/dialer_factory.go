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

var (
	errProxyJumpMissingUsername = errors.New("SSH Proxy Jump missing username")
	errProxyJumpMissingKeyFile  = errors.New("SSH Proxy Jump missing key file")
	errProxyJumpReadingKeyFile  = errors.New("SSH Proxy Jump reading key file")

	errCommandDefCreation                 = errors.New("CommandDef creation")
	errCommandDefEmptyCommand             = errors.New("command can't be empty")
	errCommandDefParsingCommandTemplate   = errors.New("parsing command template")
	errCommandDefExecutingCommandTemplate = errors.New("executing command template")
	errCommandDefInvalidCommandDefinition = errors.New("command did not contain a proper definition")
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

type DialerFactory struct {
	createDirectDialer       func(cfg Config) (*DirectDialer, error)
	createProxyJumpDialer    func(targetCfg Config, proxyCfg Config) (*ProxyJumpDialer, error)
	createProxyCommandDialer func(cfg Config, def CommandDef) (*ProxyCommandDialer, error)
	createCommandDialer      func(def CommandDef) *CommandDialer

	newCommandDef func(cmd string, target TargetDef) (CommandDef, error)
}

func NewDialerFactory() *DialerFactory {
	return &DialerFactory{
		createDirectDialer:       NewDirectDialer,
		createProxyJumpDialer:    NewProxyJumpDialer,
		createProxyCommandDialer: NewProxyCommandDialer,
		createCommandDialer:      NewCommandDialer,
		newCommandDef:            newCommandDef,
	}
}

func (df *DialerFactory) Create(flags Flags, target TargetDef) (Dialer, error) {
	if flags.Command != "" {
		return df.newCommandDialer(flags, target)
	}

	if flags.ProxyCommand != "" {
		return df.newProxyCommandDialer(flags, target)
	}

	if flags.ProxyJump.Address != "" {
		return df.newProxyJumpDialer(flags, target)
	}

	return df.newDirectDialer(flags, target)
}

func (df *DialerFactory) newCommandDialer(flags Flags, target TargetDef) (Dialer, error) {
	def, err := df.newCommandDef(flags.Command, target)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errCommandDefCreation, err)
	}

	return df.createCommandDialer(def), nil
}

func (df *DialerFactory) newProxyCommandDialer(flags Flags, target TargetDef) (Dialer, error) {
	def, err := df.newCommandDef(flags.ProxyCommand, target)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errCommandDefCreation, err)
	}

	cfg := Config{
		Address:     target.Host.Address,
		Username:    target.Host.Username,
		KeyPemBytes: target.Host.PrivateKeyPem,
	}

	return df.createProxyCommandDialer(cfg, def)
}

func (df *DialerFactory) newProxyJumpDialer(flags Flags, target TargetDef) (Dialer, error) {
	if flags.ProxyJump.Username == "" {
		return nil, errProxyJumpMissingUsername
	}

	if flags.ProxyJump.KeyFile == "" {
		return nil, errProxyJumpMissingKeyFile
	}

	keyPem, err := os.ReadFile(flags.ProxyJump.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errProxyJumpReadingKeyFile, err)
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

	return df.createProxyJumpDialer(targetCfg, proxyCfg)
}

func (df *DialerFactory) newDirectDialer(_ Flags, target TargetDef) (Dialer, error) {
	targetCfg := Config{
		Address:     target.Host.Address,
		Username:    target.Host.Username,
		KeyPemBytes: target.Host.PrivateKeyPem,
	}

	return df.createDirectDialer(targetCfg)
}

func newCommandDef(cmd string, target TargetDef) (CommandDef, error) {
	var d CommandDef

	if cmd == "" {
		return d, errCommandDefEmptyCommand
	}

	tpl, err := template.New("command").Parse(cmd)
	if err != nil {
		return d, fmt.Errorf("%w: %v", errCommandDefParsingCommandTemplate, err)
	}

	buf := bytes.NewBuffer(nil)
	err = tpl.Execute(buf, target)
	if err != nil {
		return d, fmt.Errorf("%w: %v", errCommandDefExecutingCommandTemplate, err)
	}

	evaluated := strings.TrimSpace(buf.String())
	if evaluated == "" {
		return d, errCommandDefInvalidCommandDefinition
	}

	parts := strings.Split(evaluated, " ")

	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	d.Cmd = parts[0]
	d.Args = args

	return d, nil
}
