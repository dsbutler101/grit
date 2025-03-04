package wrapper

import (
	"context"
	"fmt"
	"time"
)

const (
	CheckForRunning = true
	CheckForStopped = false

	defaultSleepTime = 1 * time.Second
)

type StatusCheckLoopTimeoutExceededError struct {
	timeout time.Duration
}

func (e *StatusCheckLoopTimeoutExceededError) Error() string {
	return fmt.Sprintf("status check loop timed out after %v", e.timeout)
}

//go:generate mockery --name=loopStatusCheckClient --inpackage --with-expecter
type loopStatusCheckClient interface {
	CheckStatus(context.Context) (Status, error)
}

func LoopStatusCheck(ctx context.Context, c loopStatusCheckClient, timeout time.Duration, checkForRunning bool) error {
	return loopStatusCheckWithSleep(ctx, c, timeout, checkForRunning, defaultSleepTime)
}

func loopStatusCheckWithSleep(ctx context.Context, c loopStatusCheckClient, timeout time.Duration, checkForRunning bool, sleep time.Duration) error {
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

		time.Sleep(sleep)
	}
}
