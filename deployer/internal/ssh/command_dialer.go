package ssh

import (
	"context"
	"errors"
	"net"
)

type CommandDef struct {
	Cmd  string
	Args []string
}

type CommandDialer struct {
	sc *streamCommand
}

func NewCommandDialer(def CommandDef) *CommandDialer {
	return &CommandDialer{
		sc: newStreamCommand(def.Cmd, def.Args...),
	}
}

func (d *CommandDialer) Start(ctx context.Context) error {
	return d.sc.Start(ctx)
}

func (d *CommandDialer) Wait() error {
	return d.sc.Wait()
}

func (d *CommandDialer) Dial(_ string, _ string) (net.Conn, error) {
	conn := d.sc.Connection()
	if conn == nil {
		return nil, errors.New("command did not provide a connection")
	}

	return conn, nil
}

func (d *CommandDialer) Close() error {
	return d.sc.Kill()
}
