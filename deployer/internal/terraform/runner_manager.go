package terraform

import (
	"errors"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
)

const (
	OutputName = "grit_runner_managers"

	OutputValueInstance       = "instance"
	OutputValueAddress        = "address"
	OutputValueWrapperAddress = "wrapper_address"
	OutputValueUsername       = "username"
	OutputValueSSHKeyPem      = "ssh_key_pem"
)

type RunnerManagers map[string]RunnerManager

type RunnerManager struct {
	Instance       string
	Address        string
	WrapperAddress string
	Username       *string
	SSHKeyPem      *string
}

func newRunnerManager(in any) (RunnerManager, error) {
	var rm RunnerManager

	rmMap, ok := in.(map[string]any)
	if !ok {
		return rm, errors.New("not a map[string]any")
	}

	type keyReader func(key string) error

	requiredStringReader := func(source map[string]any, target *string) keyReader {
		return func(key string) error {
			return readRequiredStringTo(source, key, target)
		}
	}

	optionalStringReader := func(source map[string]any, target *string) keyReader {
		return func(key string) error {
			return readOptionalStringTo(source, key, target)
		}
	}

	matchingMap := map[string]keyReader{
		OutputValueInstance:       optionalStringReader(rmMap, &rm.Instance),
		OutputValueAddress:        requiredStringReader(rmMap, &rm.Address),
		OutputValueWrapperAddress: requiredStringReader(rmMap, &rm.WrapperAddress),
		OutputValueUsername:       optionalStringReader(rmMap, rm.Username),
		OutputValueSSHKeyPem:      optionalStringReader(rmMap, rm.SSHKeyPem),
	}

	for key, reader := range matchingMap {
		err := reader(key)
		if err != nil {
			return rm, fmt.Errorf("reading output value %q: %w", key, err)
		}
	}

	return rm, nil
}

func readRunnerManagersMap(in []byte) (map[string]any, error) {
	st := new(tfjson.State)
	err := st.UnmarshalJSON(in)
	if err != nil {
		return nil, fmt.Errorf("parsing terraform state: %w", err)
	}

	if st.Values == nil {
		return nil, errors.New("terraform state does not contain any values")
	}

	if st.Values.Outputs == nil {
		return nil, errors.New("terraform state does not contain outputs")
	}

	runnerManagers, ok := st.Values.Outputs[OutputName]
	if !ok {
		return nil, fmt.Errorf("terraform state does not contain %s output", OutputName)
	}

	rmsMap, ok := runnerManagers.Value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s output is not map[string]any", OutputName)
	}

	return rmsMap, nil
}
