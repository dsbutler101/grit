package wrapper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"sync"
	"time"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api/client"
)

const (
	runnerManagerAliasLogKey = "runner-manager-alias"
)

var (
	ErrNilCallbackFunction          = errors.New("nil callback function")
	ErrWrapperAddressParsing        = errors.New("parsing wrapper address")
	ErrInvalidWrapperAddressNetwork = errors.New("invalid wrapper address network")
	ErrUsernameNotProvided          = errors.New("username not provided with CLI flag nor with Terraform Output")
	ErrSSHKeyPemNotProvided         = errors.New("ssh SSH not provided with CLI flag nor with Terraform Output")
	ErrReadingSSHKeyFile            = errors.New("reading SSH key file")
	ErrSSHDialerFabrication         = errors.New("SSH dialer fabrication")
	ErrSSHDialerStart               = errors.New("SSH dialer start")
	ErrWrapperConnect               = errors.New("wrapper connect")
	ErrWrapperExecution             = errors.New("wrapper execution")
)

type callbackFn func(ctx context.Context, c *Client) error

type runnerManagerResult struct {
	alias string
	err   error
}

//go:generate mockery --name=tfClient --inpackage --with-expecter
type tfClient interface {
	ReadStateDir(context.Context, string) (terraform.RunnerManagers, error)
	ReadStateFile(context.Context, string) (terraform.RunnerManagers, error)
}

//go:generate mockery --name=rmHandlerFactory --inpackage --with-expecter
type rmHandlerFactory interface {
	new(string) rmHandler
}

//go:generate mockery --name=rmHandler --inpackage --with-expecter
type rmHandler interface {
	handle(context.Context, terraform.RunnerManager, callbackFn) error
}

type Mux struct {
	logger *slog.Logger
	tf     tfClient

	tfFlags terraform.Flags

	rmHandlerFactory rmHandlerFactory
}

func NewMux(logger *slog.Logger, tf tfClient, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags Flags) *Mux {
	return &Mux{
		logger:  logger,
		tf:      tf,
		tfFlags: tfFlags,
		rmHandlerFactory: &muxRunnerManagerHandlerFactory{
			logger:       logger,
			sshFlags:     sshFlags,
			wrapperFlags: wrapperFlags,
		},
	}
}

func (m *Mux) Execute(ctx context.Context, fn callbackFn) error {
	if fn == nil {
		return ErrNilCallbackFunction
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

			runnerManagerHandler := m.rmHandlerFactory.new(alias)

			resultCh <- runnerManagerResult{
				alias: alias,
				err:   runnerManagerHandler.handle(ctx, runnerManager, fn),
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

type muxRunnerManagerHandlerFactory struct {
	logger *slog.Logger

	sshFlags     ssh.Flags
	wrapperFlags Flags

	dialerFactory dialerFactory
	clientFactory clientFactory
}

func (f *muxRunnerManagerHandlerFactory) new(alias string) rmHandler {
	df := ssh.NewDialer
	if f.dialerFactory != nil {
		df = f.dialerFactory
	}

	cf := NewConnectedClient
	if f.clientFactory != nil {
		cf = f.clientFactory
	}

	return &muxRunnerManagerHandler{
		logger:        f.logger.With(runnerManagerAliasLogKey, alias),
		sshFlags:      f.sshFlags,
		wrapperFlags:  f.wrapperFlags,
		dialerFactory: df,
		clientFactory: cf,
	}
}

type dialerFactory func(flags ssh.Flags, def ssh.TargetDef) (ssh.Dialer, error)

type clientFactory func(ctx context.Context, logger *slog.Logger, dialer client.Dialer, connectionTimeout time.Duration, address string) (*Client, error)

type muxRunnerManagerHandler struct {
	logger *slog.Logger

	sshFlags     ssh.Flags
	wrapperFlags Flags

	dialerFactory dialerFactory
	clientFactory clientFactory
}

func (h *muxRunnerManagerHandler) handle(ctx context.Context, rm terraform.RunnerManager, fn callbackFn) error {
	h.logger.Info("Handling runner manager")

	wa, err := getWrapperAddress(rm)
	if err != nil {
		return fmt.Errorf("getting wrapper address: %w", err)
	}

	username, err := getUsername(rm, h.sshFlags)
	if err != nil {
		return fmt.Errorf("getting username: %w", err)
	}

	sshKeyPemBytes, err := getSSHKeyPemBytes(rm, h.sshFlags)
	if err != nil {
		return fmt.Errorf("getting ssh key bytes: %w", err)
	}

	dialer, err := h.dialerFactory(h.sshFlags, ssh.TargetDef{
		Host: ssh.TargetHostDef{
			Address:       rm.Address,
			Username:      username,
			PrivateKeyPem: sshKeyPemBytes,
		},
		GRPCServer: ssh.TargetGRPCServerDef{
			Network: wa.network,
			Address: wa.address,
		},
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSSHDialerFabrication, err)
	}

	var wg sync.WaitGroup

	defer func() {
		err := dialer.Close()
		if err != nil {
			h.logger.With("error", err).Error("Closing SSH Dialer")
		}
		wg.Wait()
	}()

	err = dialer.Start(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSSHDialerStart, err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := dialer.Wait()
		if err != nil {
			h.logger.With("error", err).Error("SSH execution failure")
		}
	}()

	c, err := h.clientFactory(ctx, h.logger, dialer.Dial, h.wrapperFlags.ConnectionTimeout, rm.WrapperAddress)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWrapperConnect, err)
	}

	err = fn(ctx, c)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWrapperExecution, err)
	}

	return nil
}

type wrapperAddress struct {
	network string
	address string
}

func getWrapperAddress(rm terraform.RunnerManager) (wrapperAddress, error) {
	var wa wrapperAddress

	u, err := url.Parse(rm.WrapperAddress)
	if err != nil {
		return wa, fmt.Errorf("%w: %v", ErrWrapperAddressParsing, err)
	}

	if u.Scheme != "unix" && u.Scheme != "tcp" {
		return wa, fmt.Errorf("%w: must be unix socket or tcp socket; got %s", ErrInvalidWrapperAddressNetwork, u.Scheme)
	}

	var address string
	if u.Scheme == "unix" {
		address = u.Path
	} else {
		address = fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
	}

	wa.network = u.Scheme
	wa.address = address

	return wa, nil
}

func getUsername(rm terraform.RunnerManager, sshFlags ssh.Flags) (string, error) {
	if rm.Username == "" && sshFlags.Username == "" {
		return "", ErrUsernameNotProvided
	}

	if sshFlags.Username != "" {
		return sshFlags.Username, nil
	}

	return rm.Username, nil
}

func getSSHKeyPemBytes(rm terraform.RunnerManager, sshFlags ssh.Flags) ([]byte, error) {
	if rm.SSHKeyPem == "" && sshFlags.KeyFile == "" {
		return nil, ErrSSHKeyPemNotProvided
	}

	if sshFlags.KeyFile != "" {
		var err error
		sshKeyPemBytes, err := os.ReadFile(sshFlags.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("%w %s: %v", ErrReadingSSHKeyFile, sshFlags.KeyFile, err)
		}

		return sshKeyPemBytes, nil
	}

	return []byte(rm.SSHKeyPem), nil
}
