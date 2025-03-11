package terraform

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

type CommandError struct {
	command  string
	exitCode int
	err      error
}

func newCommandError(command string, exitCode int, err error) *CommandError {
	return &CommandError{
		command:  command,
		exitCode: exitCode,
		err:      err,
	}
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("terraform %s failed: %v", e.command, e.err)
}

func (e *CommandError) Unwrap() error {
	return e.err
}

func (e *CommandError) ExitCode() int {
	return e.exitCode
}

type Client struct {
	logger *slog.Logger

	execPath string
}

func New(logger *slog.Logger) (*Client, error) {
	execPath, err := exec.LookPath("terraform")
	if err != nil {
		return nil, fmt.Errorf("terraform not found in $PATH: %w", err)
	}

	c := &Client{
		logger:   logger,
		execPath: execPath,
	}

	return c, nil
}

func (c *Client) ExecPath() string {
	return c.execPath
}

func (c *Client) SetExecPath(execPath string) {
	c.execPath = execPath
}

func (c *Client) Init(ctx context.Context, workdir string) error {
	return c.runTerraform(ctx, workdir, "init")
}

func (c *Client) runTerraform(ctx context.Context, workdir string, command string, args ...string) error {
	return c.runTerraformWithOutput(ctx, os.Stdout, workdir, command, args...)
}

func (c *Client) runTerraformWithOutput(ctx context.Context, out io.Writer, workdir string, command string, args ...string) error {
	workdir, err := assertWorkdir(workdir)
	if err != nil {
		return err
	}

	args = append([]string{command}, args...)
	cmd := exec.CommandContext(ctx, c.execPath, args...)
	cmd.Stdout = out
	cmd.Stderr = os.Stdout
	cmd.Dir = workdir

	log := c.logger.With("command", command, "workdir", workdir)
	log.Info("Running terraform", "command", command, "workdir", workdir)

	err = cmd.Run()
	if err != nil {
		log.Error("Terraform execution failed", "command", command, "error", err, "exit-code", cmd.ProcessState.ExitCode())

		return newCommandError(command, cmd.ProcessState.ExitCode(), err)
	}

	log.Info("Terraform execution succeeded", "command", command)

	return nil
}

func assertWorkdir(wd string) (string, error) {
	if filepath.IsAbs(wd) {
		return wd, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("requesting current working directory: %w", err)
	}

	return filepath.Join(dir, wd), nil
}

func (c *Client) Apply(ctx context.Context, workdir string) error {
	return c.runTerraform(ctx, workdir, "apply", "--auto-approve")
}

func (c *Client) ReadStateDir(ctx context.Context, workdir string) (RunnerManagers, error) {
	buf := new(bytes.Buffer)
	err := c.runTerraformWithOutput(ctx, buf, workdir, "show", "-json")
	if err != nil {
		return RunnerManagers{}, fmt.Errorf("reading terraform state: %w", err)
	}

	return c.readState(buf)
}

func (c *Client) ReadStateFile(_ context.Context, filePath string) (RunnerManagers, error) {
	buf := new(bytes.Buffer)
	f, err := os.Open(filePath)
	if err != nil {
		return RunnerManagers{}, fmt.Errorf("opening terraform state file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(buf, f)
	if err != nil {
		return RunnerManagers{}, fmt.Errorf("reading terraform state file: %w", err)
	}

	return c.readState(buf)
}

func (c *Client) readState(buf *bytes.Buffer) (RunnerManagers, error) {
	var output RunnerManagers

	rmsMap, err := readRunnerManagersMap(buf.Bytes())
	if err != nil {
		return output, fmt.Errorf("reading terraform state: %w", err)
	}

	output = make(RunnerManagers, len(rmsMap))

	for name, runnerManager := range rmsMap {
		rm, err := newRunnerManager(runnerManager)
		if err != nil {
			return output, fmt.Errorf("reading runner manager %q: %w", name, err)
		}

		output[name] = rm
	}

	return output, nil
}

func (c *Client) Destroy(ctx context.Context, workdir string) error {
	return c.runTerraform(ctx, workdir, "destroy", "--auto-approve")
}
