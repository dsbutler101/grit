package wrapper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v3"
	crypto_ssh "golang.org/x/crypto/ssh"

	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api/client"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/logger"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/ssh"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/deployer/internal/terraform"
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

type CallbackFn func(ctx context.Context, c CallbackClient) error

//go:generate mockery --name=CallbackClient --inpackage --with-expecter
type CallbackClient interface {
	CheckStatus(context.Context) (Status, error)
	InitForcefulShutdown(context.Context) error
	InitGracefulShutdown(context.Context) error
}

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
	handle(context.Context, terraform.RunnerManager, CallbackFn) error
}

//go:generate mockery --name=dialerFactory --inpackage --with-expecter
type dialerFactory interface {
	Create(flags ssh.Flags, def ssh.TargetDef) (ssh.Dialer, error)
}

type Mux struct {
	logger *slog.Logger
	tf     tfClient

	tfFlags      terraform.Flags
	wrapperFlags Flags

	rmHandlerFactory rmHandlerFactory
}

func NewMux(logger *slog.Logger, tf tfClient, tfFlags terraform.Flags, sshFlags ssh.Flags, wrapperFlags Flags) *Mux {
	return &Mux{
		logger:       logger,
		tf:           tf,
		tfFlags:      tfFlags,
		wrapperFlags: wrapperFlags,
		rmHandlerFactory: &muxRunnerManagerHandlerFactory{
			logger:       logger,
			sshFlags:     sshFlags,
			wrapperFlags: wrapperFlags,
		},
	}
}

func (m *Mux) Execute(ctx context.Context, fn CallbackFn) error {
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

			err := m.handleRunnerManager(ctx, alias, runnerManager, fn, m.wrapperFlags.ConnectionTimeout)

			resultCh <- runnerManagerResult{
				alias: alias,
				err:   err,
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
				log.With(logger.ErrorKey, result.err).Error("Runner manager execution failed")
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

func (m *Mux) handleRunnerManager(ctx context.Context, alias string, rm terraform.RunnerManager, fn CallbackFn, connectionTimeout time.Duration) error {
	cctx, cancelFn := context.WithTimeout(ctx, connectionTimeout)
	defer cancelFn()

	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = 5 * time.Second

	for {
		select {
		case <-cctx.Done():
			err := ctx.Err()
			if err != nil {
				return err
			}

			return client.ErrRetryTimeoutExceeded
		default:
			err := m.rmHandlerFactory.new(alias).handle(ctx, rm, fn)

			var eerr *net.OpError
			if errors.As(err, &eerr) {
				sleep := bo.NextBackOff()
				m.logger.With(logger.ErrorKey, eerr, "sleep_seconds", sleep.Seconds()).Warn("Connection failure; retrying")

				time.Sleep(sleep)

				continue
			}

			return err
		}
	}
}

type muxRunnerManagerHandlerFactory struct {
	logger *slog.Logger

	sshFlags     ssh.Flags
	wrapperFlags Flags

	dialerFactory dialerFactory
	clientFactory clientFactory
}

func (f *muxRunnerManagerHandlerFactory) new(alias string) rmHandler {
	var df dialerFactory = ssh.NewDialerFactory()
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

type clientFactory func(ctx context.Context, logger *slog.Logger, dialer client.Dialer, connectionTimeout time.Duration, address string) (*Client, error)

type muxRunnerManagerHandler struct {
	logger *slog.Logger

	sshFlags     ssh.Flags
	wrapperFlags Flags

	dialerFactory dialerFactory
	clientFactory clientFactory
}

func (h *muxRunnerManagerHandler) handle(ctx context.Context, rm terraform.RunnerManager, fn CallbackFn) error {
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

	var fingerprintSHA256 string
	var fingerprintMD5 string
	var keyType string

	if len(sshKeyPemBytes) > 0 {
		k, err := crypto_ssh.ParsePrivateKey(sshKeyPemBytes)
		if err == nil {

			fingerprintSHA256 = crypto_ssh.FingerprintSHA256(k.PublicKey())
			fingerprintMD5 = crypto_ssh.FingerprintLegacyMD5(k.PublicKey())
			keyType = k.PublicKey().Type()
		}
	}

	log := h.logger.
		WithGroup("dialer").
		With("ssh-config", struct {
			Address              string
			Username             string
			KeyFingerprintSHA256 string
			KeyFingerprintMD5    string
			KeyType              string
		}{
			Address:              rm.Address,
			Username:             username,
			KeyFingerprintSHA256: fingerprintSHA256,
			KeyFingerprintMD5:    fingerprintMD5,
			KeyType:              keyType,
		}).
		With("grpc-config", struct {
			Network string
			Address string
		}{
			Network: wa.network,
			Address: wa.address,
		})
	log.Info("Creating SSH dialer")

	dialer, err := h.dialerFactory.Create(h.sshFlags, ssh.TargetDef{
		Host: ssh.TargetHostDef{
			InstanceName:  rm.InstanceName,
			InstanceID:    rm.InstanceID,
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
		return fmt.Errorf("%w: %w", ErrSSHDialerFabrication, err)
	}

	var wg sync.WaitGroup

	defer func() {
		err := dialer.Close()
		if err != nil {
			h.logger.With(logger.ErrorKey, err).Error("Closing SSH Dialer")
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
			h.logger.With(logger.ErrorKey, err).Error("SSH execution failure")
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
