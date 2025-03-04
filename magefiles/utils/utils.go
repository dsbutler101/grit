package utils

import (
	"fmt"
	"os"
	"strconv"
)

const (
	ciEnv = "CI"
)

func OnWD(path string, fn func() error) error {
	origWd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	err = os.Chdir(path)
	if err != nil {
		return fmt.Errorf("failed to change directory to %s: %w", origWd, err)
	}

	defer func() {
		err := os.Chdir(origWd)
		if err != nil {
			panic(fmt.Sprintf("failed to restore original working directory: %v", err))
		}
	}()

	return fn()
}

func IsCI() bool {
	ci := os.Getenv(ciEnv)
	if ci == "" {
		return false
	}

	ok, err := strconv.ParseBool(ci)
	if err != nil {
		return false
	}

	return ok
}
