package wrapper

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api"
	"gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api/client"
)

const (
	DefaultTimeout = 30 * time.Second
)

type Status struct {
	Status        api.Status
	FailureReason string
}

func (s Status) IsRunning() bool {
	return s.Status == api.StatusRunning || s.Status == api.StatusInShutdown
}

type Client struct {
	logger *slog.Logger

	c *client.Client
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
		return &GRPCConnectionRetryExceededError{
			timeout: connectionTimeout,
			err:     err,
		}
	}

	return nil
}

func (c *Client) CheckStatus(ctx context.Context) (Status, error) {
	s, err := c.c.CheckStatus(ctx)
	if err != nil {
		return Status{}, err
	}

	status := Status{
		Status:        s.Status,
		FailureReason: s.FailureReason,
	}

	return status, nil
}

func (c *Client) InitGracefulShutdown(ctx context.Context) error {
	s, err := c.c.InitGracefulShutdown(ctx, api.NewInitGracefulShutdownRequest(nil))
	if err != nil {
		return err
	}

	c.logger.WithGroup("response").With("status", s.Status.String(), "failureReason", s.FailureReason).Info("Graceful shutdown started")

	return nil
}

func (c *Client) InitForcefulShutdown(ctx context.Context) error {
	s, err := c.c.InitForcefulShutdown(ctx)
	if err != nil {
		return err
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
