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

type multiReaderCloserCloseErr struct {
	failures uint
	lastErr  error
}

func (e *multiReaderCloserCloseErr) Error() string {
	return fmt.Sprintf("encountered %d error when closing multi reader sources; last was: %v", e.failures, e.lastErr)
}

type multiReadCloser struct {
	sources []io.ReadCloser
	r       io.Reader
}

func newMultiReadCloser(readers ...io.ReadCloser) io.ReadCloser {
	r := make([]io.Reader, len(readers))
	for no, in := range readers {
		r[no] = in
	}

	return &multiReadCloser{
		r:       io.MultiReader(r...),
		sources: readers,
	}
}

func (m *multiReadCloser) Read(p []byte) (n int, err error) {
	return m.r.Read(p)
}

func (m *multiReadCloser) Close() error {
	var err error

	failures := uint(0)
	for _, source := range m.sources {
		err = source.Close()
		if err != nil {
			failures++
		}
	}

	if err != nil {
		return &multiReaderCloserCloseErr{
			failures: failures,
			lastErr:  err,
		}
	}

	return nil
}

type streamCommand struct {
	command string
	args    []string

	conn *streamConn
	mux  sync.Mutex

	cmd *exec.Cmd
}

func newStreamCommand(cmd string, args ...string) *streamCommand {
	return &streamCommand{
		command: cmd,
		args:    args,
	}
}

func (c *streamCommand) Start(ctx context.Context) error {
	if c.conn != nil || c.cmd != nil {
		return fmt.Errorf("proxy command is already running")
	}

	c.cmd = exec.CommandContext(ctx, c.command, c.args...)

	outR, err := c.cmd.StdoutPipe()
	if err != nil {
		c.cmd = nil
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	errR, err := c.cmd.StderrPipe()
	if err != nil {
		c.cmd = nil
		return fmt.Errorf("creating stderr pipe: %w", err)
	}

	inW, err := c.cmd.StdinPipe()
	if err != nil {
		c.cmd = nil
		return fmt.Errorf("creating stdin pipe: %w", err)
	}

	c.conn = newStreamConn(inW, newMultiReadCloser(outR, errR))

	return c.cmd.Start()
}

func (c *streamCommand) Wait() error {
	if c.cmd == nil {
		return fmt.Errorf("proxy command is not started")
	}

	err := c.cmd.Wait()
	if err != nil {
		var eerr *exec.ExitError
		if errors.As(err, &eerr) {
			status := eerr.ProcessState.Sys().(syscall.WaitStatus)
			if !status.Signaled() && !status.Stopped() {
				return err
			}
		}
	}

	return nil
}

func (c *streamCommand) Connection() net.Conn {
	return c.conn
}

func (c *streamCommand) Kill() error {
	if c.cmd == nil || c.cmd.Process == nil {
		return nil
	}

	return c.cmd.Process.Kill()
}
