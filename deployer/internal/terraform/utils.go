package terraform

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	errIsNotAString         = errors.New("is not a string")
	errMissingRequiredEntry = errors.New("missing required entry")
	errRequiredEntryIsEmpty = errors.New("required entry is empty")
)

func readRequiredStringTo(rmMap map[string]any, key string, target *string) error {
	return readStringTo(rmMap, key, target, false)
}

func readOptionalStringTo(rmMap map[string]any, key string, target *string) error {
	return readStringTo(rmMap, key, target, true)
}

func readStringTo(rmMap map[string]any, key string, target *string, optional bool) error {
	src, ok := rmMap[key]
	if !ok {
		if optional {
			return nil
		}

		return fmt.Errorf("%w: %s", errMissingRequiredEntry, key)
	}

	value, ok := src.(string)
	if !ok {
		return fmt.Errorf("%s %w; got %T", key, errIsNotAString, src)
	}

	if value == "" && !optional {
		return fmt.Errorf("%w: %s", errRequiredEntryIsEmpty, key)
	}

	*target = value

	return nil
}

func assertWorkdir(wd string) (string, error) {
	if filepath.IsAbs(wd) {
		return wd, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("requesting current working directory: %w", err)
	}

	return filepath.Join(dir, wd), nil
}
