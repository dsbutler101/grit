package wrapper

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api/client"
)

const (
	DefaultTimeout = 30 * time.Second
)

var (
	ErrCheckStatus          = errors.New("check status request failed")
	ErrInitGracefulShutdown = errors.New("graceful shutdown initialization failed")
	ErrInitForcefulShutdown = errors.New("forceful shutdown initialization failed")
)

type GRPCConnectionWaitTimeoutExceededError struct {
	timeout time.Duration
	err     error
}

func (e *GRPCConnectionWaitTimeoutExceededError) Error() string {
	return fmt.Sprintf("GRPC connection timeout %s exceeded: %v", e.timeout, e.err)
}

func (e *GRPCConnectionWaitTimeoutExceededError) Unwrap() error {
	return e.err
}

//go:generate mockery --name=grpcClient --inpackage --with-expecter
type grpcClient interface {
	ConnectWithTimeout(context.Context, time.Duration) error
	CheckStatus(context.Context) (client.CheckStatusResponse, error)
	InitGracefulShutdown(context.Context, api.InitGracefulShutdownRequest) (client.CheckStatusResponse, error)
	InitForcefulShutdown(context.Context) (client.CheckStatusResponse, error)
}

type Client struct {
	logger *slog.Logger

	c grpcClient
}

func NewClient(logger *slog.Logger, dialer client.Dialer, address string) (*Client, error) {
	opts := []client.Option{
		client.WithLogger(logger),
	}

	if dialer != nil {
		opts = append(opts, client.WithDialer(dialer))
	}

	c, err := client.New(address, opts...)
	if err != nil {
		return nil, err
	}

	cl := &Client{
		logger: logger,
		c:      c,
	}

	return cl, nil
}

func (c *Client) Connect(ctx context.Context, connectionTimeout time.Duration) error {
	err := c.c.ConnectWithTimeout(ctx, connectionTimeout)
	if err != nil {
		return &GRPCConnectionWaitTimeoutExceededError{
			timeout: connectionTimeout,
			err:     err,
		}
	}

	return nil
}

func (c *Client) CheckStatus(ctx context.Context) (Status, error) {
	s, err := c.c.CheckStatus(ctx)
	if err != nil {
		return Status{}, fmt.Errorf("%w: %v", ErrCheckStatus, err)
	}

	status := Status{
		Status:        s.Status,
		FailureReason: s.FailureReason,
	}

	c.logger.WithGroup("status").With("status", s.Status.String(), "failureReason", s.FailureReason).Info("Status information received")

	return status, nil
}

func (c *Client) InitGracefulShutdown(ctx context.Context) error {
	s, err := c.c.InitGracefulShutdown(ctx, api.NewInitGracefulShutdownRequest(nil))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInitGracefulShutdown, err)
	}

	c.logger.WithGroup("response").With("status", s.Status.String(), "failureReason", s.FailureReason).Info("Graceful shutdown started")

	return nil
}

func (c *Client) InitForcefulShutdown(ctx context.Context) error {
	s, err := c.c.InitForcefulShutdown(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInitForcefulShutdown, err)
	}

	c.logger.WithGroup("response").With("status", s.Status.String(), "failureReason", s.FailureReason).Info("Forceful shutdown started")

	return nil
}

func NewConnectedClient(ctx context.Context, logger *slog.Logger, dialer client.Dialer, connectionTimeout time.Duration, address string) (*Client, error) {
	c, err := NewClient(logger, dialer, address)
	if err != nil {
		return nil, fmt.Errorf("creating wrapper client: %w", err)
	}

	err = c.Connect(ctx, connectionTimeout)
	if err != nil {
		return nil, fmt.Errorf("connecting to gRPC server: %w", err)
	}

	return c, nil
}
