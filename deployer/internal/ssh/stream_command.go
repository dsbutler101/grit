package ssh

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os/exec"
	"sync"
	"syscall"
)

var (
	errStreamCommandAlreadyRunning     = errors.New("stream command is already running")
	errStreamCommandCreatingStdoutPipe = errors.New("creating stdout pipe")
	errStreamCommandCreatingStderrPipe = errors.New("creating stderr pipe")
	errStreamCommandCreatingStdinPipe  = errors.New("creating stdin pipe")
	errStreamCommandNotStarted         = errors.New("stream command is not started")
)

//go:generate mockery --name=execCmd --inpackage --with-expecter
type execCmd interface {
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
	StdinPipe() (io.WriteCloser, error)
	Start() error
	Wait() error
	Kill() error
}

type defaultExecCmd struct {
	*exec.Cmd
}

func (c *defaultExecCmd) Kill() error {
	if c.Cmd == nil || c.Cmd.Process == nil {
		return nil
	}

	return c.Cmd.Process.Kill()
}

type streamCommand struct {
	command string
	args    []string

	commandFactory func(ctx context.Context, command string, args ...string) execCmd

	conn *streamConn
	mux  sync.Mutex

	cmd execCmd
}

func newStreamCommand(cmd string, args ...string) *streamCommand {
	return &streamCommand{
		command: cmd,
		args:    args,
		commandFactory: func(ctx context.Context, command string, args ...string) execCmd {
			return &defaultExecCmd{Cmd: exec.CommandContext(ctx, command, args...)}
		},
	}
}

func (c *streamCommand) Start(ctx context.Context) error {
	if c.conn != nil || c.cmd != nil {
		return errStreamCommandAlreadyRunning
	}

	c.cmd = c.commandFactory(ctx, c.command, c.args...)

	outR, err := c.cmd.StdoutPipe()
	if err != nil {
		c.cmd = nil
		return fmt.Errorf("%w: %v", errStreamCommandCreatingStdoutPipe, err)
	}

	errR, err := c.cmd.StderrPipe()
	if err != nil {
		c.cmd = nil
		return fmt.Errorf("%w: %v", errStreamCommandCreatingStderrPipe, err)
	}

	inW, err := c.cmd.StdinPipe()
	if err != nil {
		c.cmd = nil
		return fmt.Errorf("%w: %v", errStreamCommandCreatingStdinPipe, err)
	}

	c.conn = newStreamConn(inW, newMultiReadCloser(outR, errR))

	return c.cmd.Start()
}

func (c *streamCommand) Wait() error {
	if c.cmd == nil {
		return errStreamCommandNotStarted
	}

	err := c.cmd.Wait()
	if err != nil {
		var eerr *exec.ExitError
		if errors.As(err, &eerr) {
			status := eerr.ProcessState.Sys().(syscall.WaitStatus)
			if !status.Signaled() && !status.Stopped() {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (c *streamCommand) Connection() net.Conn {
	return c.conn
}

func (c *streamCommand) Kill() error {
	return c.cmd.Kill()
}
