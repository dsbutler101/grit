package terraform

import (
	"context"
	"io"
	"os/exec"
)

const (
	CommanderCommandNotExecutedExitCode = -100
)

//go:generate mockery --name=commander --inpackage --with-expecter
type commander interface {
	run(ctx context.Context, cmd string, args ...string) error
	exitCode() int
	setStdout(io.Writer)
}

type defaultCommander struct {
	stdout  io.Writer
	stderr  io.Writer
	workdir string

	cmd *exec.Cmd
}

func newDefaultCommander(stdout io.Writer, stderr io.Writer, workdir string) commander {
	cmdr := &defaultCommander{
		stderr:  stderr,
		workdir: workdir,
	}
	cmdr.setStdout(stdout)

	return cmdr
}

func (d *defaultCommander) run(ctx context.Context, cmd string, args ...string) error {
	d.cmd = exec.CommandContext(ctx, cmd, args...)
	d.cmd.Stdout = d.stdout
	d.cmd.Stderr = d.stderr
	d.cmd.Dir = d.workdir

	return d.cmd.Run()
}

func (d *defaultCommander) setStdout(stdout io.Writer) {
	d.stdout = stdout
}

func (d *defaultCommander) exitCode() int {
	if d.cmd == nil {
		return CommanderCommandNotExecutedExitCode
	}

	return d.cmd.ProcessState.ExitCode()
}
