package wrapper

import (
	"fmt"
	"time"
)

type GRPCConnectionRetryExceededError struct {
	timeout time.Duration
	err     error
}

func (e *GRPCConnectionRetryExceededError) Error() string {
	return fmt.Sprintf("GRPC connection timeout %s exceeded: %v", e.timeout, e.err)
}

func (e *GRPCConnectionRetryExceededError) Unwrap() error {
	return e.err
}
