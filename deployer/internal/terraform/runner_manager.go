package terraform

import (
	"errors"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
)

const (
	OutputName = "grit_runner_managers"

	OutputValueInstanceName   = "instance_name"
	OutputValueAddress        = "address"
	OutputValueWrapperAddress = "wrapper_address"
	OutputValueUsername       = "username"
	OutputValueSSHKeyPem      = "ssh_key_pem"
)

var (
	errIsNotAMapString = errors.New("is not a map[string]any")

	errParsingTerraformState               = errors.New("parsing terraform state")
	errTerraformStateMissingValues         = errors.New("terraform state does not contain any values")
	errTerraformStateMissingOutputs        = errors.New("terraform state does not contain outputs")
	errTerraformStateMissingRequiredOutput = errors.New("terraform state does not contain required output")
)

type RunnerManagers map[string]RunnerManager

type RunnerManager struct {
	InstanceName   string
	Address        string
	WrapperAddress string
	Username       string
	SSHKeyPem      string
}

func newRunnerManager(in any) (RunnerManager, error) {
	var rm RunnerManager

	rmMap, ok := in.(map[string]any)
	if !ok {
		return rm, errIsNotAMapString
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
		OutputValueInstanceName:   optionalStringReader(rmMap, &rm.InstanceName),
		OutputValueAddress:        requiredStringReader(rmMap, &rm.Address),
		OutputValueWrapperAddress: requiredStringReader(rmMap, &rm.WrapperAddress),
		OutputValueUsername:       optionalStringReader(rmMap, &rm.Username),
		OutputValueSSHKeyPem:      optionalStringReader(rmMap, &rm.SSHKeyPem),
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
		return nil, fmt.Errorf("%w: %v", errParsingTerraformState, err)
	}

	if st.Values == nil {
		return nil, errTerraformStateMissingValues
	}

	if st.Values.Outputs == nil {
		return nil, errTerraformStateMissingOutputs
	}

	runnerManagers, ok := st.Values.Outputs[OutputName]
	if !ok {
		return nil, fmt.Errorf("%w: %s", errTerraformStateMissingRequiredOutput, OutputName)
	}

	rmsMap, ok := runnerManagers.Value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s %w", OutputName, errIsNotAMapString)
	}

	return rmsMap, nil
}
