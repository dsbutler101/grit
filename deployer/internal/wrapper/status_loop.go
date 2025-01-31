package wrapper

import (
	"context"
	"fmt"
	"time"
)

type StatusCheckLoopTimeoutExceededError struct {
	timeout time.Duration
}

func (e *StatusCheckLoopTimeoutExceededError) Error() string {
	return fmt.Sprintf("status check loop timed out after %v", e.timeout)
}

func LoopStatusCheck(ctx context.Context, c *Client, timeout time.Duration, checkForRunning bool) error {
	startTime := time.Now()
	for {
		status, err := c.CheckStatus(ctx)
		if err != nil {
			return fmt.Errorf("checking status: %w", err)
		}

		if status.IsRunning() == checkForRunning {
			return nil
		}

		if time.Now().Sub(startTime) > timeout {
			return &StatusCheckLoopTimeoutExceededError{timeout}
		}

		time.Sleep(1 * time.Second)
	}
}
