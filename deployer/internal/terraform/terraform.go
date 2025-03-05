package terraform

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
)

const (
	tfBinary         = "terraform"
	tfCommandInit    = "init"
	tfCommandApply   = "apply"
	tfCommandDestroy = "destroy"
	tfCommandShow    = "show"

	tfAutoApproveFlag = "--auto-approve"
)

var (
	errReadingTerraformState = errors.New("reading terraform state")
	errReadingRunnerManager  = errors.New("reading runner manager")

	errReadingTerraformStateFile = errors.New("reading terraform state file")
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

var (
	ErrTerraformNotFoundInPath = errors.New("terraform not found in $PATH")
	ErrTerraformExecPathNotSet = errors.New("terraform exec path is not set")
)

func DefaultExecPath() (string, error) {
	execPath, err := exec.LookPath(tfBinary)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTerraformNotFoundInPath, err)
	}

	return execPath, nil
}

type stateReader func(buf *bytes.Buffer) (RunnerManagers, error)

type Client struct {
	logger *slog.Logger

	commanderFactory func(stdout io.Writer, stderr io.Writer, workdir string) commander
	stateReader      stateReader

	execPath string
}

func New(logger *slog.Logger) *Client {
	c := &Client{
		logger:           logger,
		commanderFactory: newDefaultCommander,
		stateReader:      readState,
	}

	return c
}

func (c *Client) SetExecPath(execPath string) {
	c.execPath = execPath
}

func (c *Client) Init(ctx context.Context, workdir string) error {
	return c.runTerraform(ctx, workdir, tfCommandInit)
}

func (c *Client) Apply(ctx context.Context, workdir string) error {
	return c.runTerraform(ctx, workdir, tfCommandApply, tfAutoApproveFlag)
}

func (c *Client) Destroy(ctx context.Context, workdir string) error {
	return c.runTerraform(ctx, workdir, tfCommandDestroy, tfAutoApproveFlag)
}

func (c *Client) runTerraform(ctx context.Context, workdir string, command string, args ...string) error {
	return c.runTerraformWithOutput(ctx, os.Stdout, workdir, command, args...)
}

func (c *Client) runTerraformWithOutput(ctx context.Context, out io.Writer, workdir string, command string, args ...string) error {
	if c.execPath == "" {
		return ErrTerraformExecPathNotSet
	}

	workdir, err := assertWorkdir(workdir)
	if err != nil {
		return err
	}

	log := c.logger.With("command", command, "workdir", workdir)
	log.Info("Running terraform")

	args = append([]string{command}, args...)
	cmd := c.commanderFactory(out, os.Stdout, workdir)

	err = cmd.run(ctx, c.execPath, args...)
	if err != nil {
		log.Error("Terraform execution failed", "error", err, "exit-code", cmd.exitCode())

		return newCommandError(command, cmd.exitCode(), err)
	}

	log.Info("Terraform execution succeeded", "command", command)

	return nil
}

func (c *Client) ReadStateDir(ctx context.Context, workdir string) (RunnerManagers, error) {
	buf := new(bytes.Buffer)
	err := c.runTerraformWithOutput(ctx, buf, workdir, tfCommandShow, "-json")
	if err != nil {
		return RunnerManagers{}, fmt.Errorf("reading terraform state: %w", err)
	}

	return c.stateReader(buf)
}

func (c *Client) ReadStateFile(_ context.Context, filePath string) (RunnerManagers, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return RunnerManagers{}, fmt.Errorf("%w: %w", errReadingTerraformStateFile, err)
	}

	return c.stateReader(bytes.NewBuffer(data))
}

func readState(buf *bytes.Buffer) (RunnerManagers, error) {
	var output RunnerManagers

	rmsMap, err := readRunnerManagersMap(buf.Bytes())
	if err != nil {
		return output, fmt.Errorf("%w: %w", errReadingTerraformState, err)
	}

	output = make(RunnerManagers, len(rmsMap))

	for name, runnerManager := range rmsMap {
		rm, err := newRunnerManager(runnerManager)
		if err != nil {
			return output, fmt.Errorf("%w %q: %w", errReadingRunnerManager, name, err)
		}

		output[name] = rm
	}

	return output, nil
}
