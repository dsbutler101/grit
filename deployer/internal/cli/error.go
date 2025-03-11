package cli

import (
	"fmt"
)

type Error struct {
	exitCode int
	err      error
}

func NewError(exitCode int, err error) *Error {
	return &Error{
		exitCode: exitCode,
		err:      err,
	}
}

func (e *Error) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	return fmt.Sprintf("unknown error (exit code %d)", e.exitCode)
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) ExitCode() int {
	return e.exitCode
}
