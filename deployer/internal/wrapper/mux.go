package wrapper

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"sync"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
)

const (
	runnerManagerAliasLogKey = "runner-manager-alias"
)

type callbackFn func(ctx context.Context, c *Client) error

type runnerManagerResult struct {
	alias string
	err   error
}

type Mux struct {
	logger *slog.Logger
	tf     *terraform.Client

	tfFlags      terraform.Flags
	sshFlags     ssh.Flags
	wrapperFlags Flags
}

func NewMux(logger *slog.Logger, tf *terraform.Client, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags Flags) *Mux {
	return &Mux{
		logger:       logger,
		tf:           tf,
		tfFlags:      tfFlags,
		sshFlags:     sshFlags,
		wrapperFlags: wrapperFlags,
	}
}

func (m *Mux) Execute(ctx context.Context, fn callbackFn) error {
	if fn == nil {
		return fmt.Errorf("nil callback function")
	}

	runnerManagers, err := m.listRunnerManagers(ctx)
	if err != nil {
		return err
	}

	resultCh := make(chan runnerManagerResult, len(runnerManagers))
	defer close(resultCh)

	finishedCh := make(chan struct{})

	wg := new(sync.WaitGroup)
	for alias, runnerManager := range runnerManagers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			resultCh <- runnerManagerResult{
				alias: alias,
				err:   m.handleRunnerManager(ctx, runnerManager, alias, fn),
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finishedCh)
	}()

	var lastErr error
	for {
		select {
		case <-ctx.Done():
			return lastErr
		case result := <-resultCh:
			log := m.logger.With(runnerManagerAliasLogKey, result.alias)
			if result.err != nil {
				log.With("error", result.err).Error("Runner manager execution failed")
				lastErr = result.err
				continue
			}

			log.Info("Runner manager execution succeeded")
		case <-finishedCh:
			return lastErr
		}
	}
}

func (m *Mux) listRunnerManagers(ctx context.Context) (terraform.RunnerManagers, error) {
	if m.tfFlags.Target != "" {
		return m.tf.ReadStateDir(ctx, m.tfFlags.Target)
	}

	return m.tf.ReadStateFile(ctx, m.tfFlags.TargetStateFile)
}

func (m *Mux) handleRunnerManager(ctx context.Context, rm terraform.RunnerManager, alias string, fn callbackFn) error {
	log := m.logger.With("runner-manager", rm, runnerManagerAliasLogKey, alias)
	log.Info("Handling runner manager")

	wa, err := m.getWrapperAddress(rm)
	if err != nil {
		return fmt.Errorf("getting wrapper address: %w", err)
	}

	username, err := m.getUsername(rm)
	if err != nil {
		return fmt.Errorf("getting username: %w", err)
	}

	sshKeyPemBytes, err := m.getSSHKeyPemBytes(rm)
	if err != nil {
		return fmt.Errorf("getting ssh key bytes: %w", err)
	}

	dialer, err := ssh.NewDialer(m.sshFlags, ssh.TargetDef{
		Host: ssh.TargetHostDef{
			Address:       rm.Address,
			Username:      username,
			PrivateKeyPem: sshKeyPemBytes,
		},
		GRPCServer: ssh.TargetGRPCServerDef{
			Network: wa.scheme,
			Address: wa.address,
		},
	})
	if err != nil {
		return fmt.Errorf("SSH Dialer fabrication: %w", err)
	}

	defer func() {
		err := dialer.Close()
		if err != nil {
			log.With("error", err).Error("Closing SSH Dialer")
		}
	}()

	err = dialer.Start(ctx)
	if err != nil {
		return fmt.Errorf("SSH Dialer start: %w", err)
	}

	go func() {
		err := dialer.Wait()
		if err != nil {
			log.With("error", err).Error("SSH execution failure")
		}
	}()

	c, err := NewConnectedClient(ctx, log, dialer.Dial, m.wrapperFlags.ConnectionTimeout, rm.WrapperAddress)
	if err != nil {
		return fmt.Errorf("wrapper connect: %w", err)
	}

	err = fn(ctx, c)
	if err != nil {
		return fmt.Errorf("wrapper execution: %w", err)
	}

	return nil
}

type wrapperAddress struct {
	scheme  string
	address string
}

func (m *Mux) getWrapperAddress(rm terraform.RunnerManager) (wrapperAddress, error) {
	var wa wrapperAddress

	u, err := url.Parse(rm.WrapperAddress)
	if err != nil {
		return wa, fmt.Errorf("parsing wrapper address: %w", err)
	}

	if u.Scheme != "unix" && u.Scheme != "tcp" {
		return wa, fmt.Errorf("wrapper address must be unix socket or tcp socket; got %s", u.Scheme)
	}

	var address string
	if u.Scheme == "unix" {
		address = u.Path
	} else {
		address = fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
	}

	wa.scheme = u.Scheme
	wa.address = address

	return wa, nil
}

func (m *Mux) getUsername(rm terraform.RunnerManager) (string, error) {
	if rm.Username == nil && m.sshFlags.Username == "" {
		return "", fmt.Errorf("username not provided with CLI flag nor with Terraform Output")
	}

	if m.sshFlags.Username != "" {
		return m.sshFlags.Username, nil
	}

	return *rm.Username, nil
}

func (m *Mux) getSSHKeyPemBytes(rm terraform.RunnerManager) ([]byte, error) {
	var sshKeyPemBytes []byte
	if rm.SSHKeyPem != nil {
		sshKeyPemBytes = []byte(*rm.SSHKeyPem)
	}

	if m.sshFlags.KeyFile != "" {
		var err error
		sshKeyPemBytes, err = os.ReadFile(m.sshFlags.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("reading ssh key file %q: %w", m.sshFlags.KeyFile, err)
		}
	}

	return sshKeyPemBytes, nil
}
