package terraform

import (
	"fmt"
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

		return fmt.Errorf("does not contain %q output", key)
	}

	value, ok := src.(string)
	if !ok {
		return fmt.Errorf("%s is not a string; got %T", key, src)
	}

	*target = value

	return nil
}
