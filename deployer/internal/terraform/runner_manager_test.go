package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRunnerManager(t *testing.T) {
	testInstanceName := "test-instance-name"
	testAddress := "localhost:1234"
	testWrapperAddress := "unix:///var/run/wrapper.sock"

	tests := map[string]struct {
		input       any
		assertError func(t *testing.T, err error)
		assertRM    func(t *testing.T, rm RunnerManager)
	}{
		"invalid input type": {
			input: "invalid",
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errIsNotAMapString)
			},
		},
		"input reading error": {
			input: map[string]any{
				"instance_name": testInstanceName,
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errMissingRequiredEntry)
			},
		},
		"correct input": {
			input: map[string]any{
				"instance_name":   testInstanceName,
				"address":         testAddress,
				"wrapper_address": testWrapperAddress,
			},
			assertRM: func(t *testing.T, rm RunnerManager) {
				assert.Equal(t, testInstanceName, rm.InstanceName)
				assert.Equal(t, testAddress, rm.Address)
				assert.Equal(t, testWrapperAddress, rm.WrapperAddress)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			rm, err := newRunnerManager(tt.input)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, tt.assertRM, "assertRM must be defined in test definition")
			tt.assertRM(t, rm)
		})
	}
}

func TestReadRunnerManagersMap(t *testing.T) {
	tests := map[string]struct {
		input       string
		assertError func(t *testing.T, err error)
		assertMap   func(t *testing.T, m map[string]any)
	}{
		"not a valid Terraform State JSON": {
			input: ``,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errParsingTerraformState)
			},
		},
		"values missing in Terraform State": {
			input: `{"format_version": "1"}`,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errTerraformStateMissingValues)
			},
		},
		"outputs missing in Terraform State": {
			input: `{"format_version": "1", "values": {}}`,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errTerraformStateMissingOutputs)
			},
		},
		"grit_runner_managers missing in Terraform State outputs": {
			input: `{"format_version": "1", "values": {"outputs": {}}}`,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errTerraformStateMissingRequiredOutput)
			},
		},
		"grit_runner_managers Terraform State output is not a map[string]any": {
			input: `{"format_version": "1", "values": {"outputs": {"grit_runner_managers": {"value": "test"}}}}`,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errIsNotAMapString)
			},
		},
		"grit_runner_managers Terraform State output is null": {
			input: `{"format_version": "1", "values": {"outputs": {"grit_runner_managers": {"value": null}}}}`,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, errIsNotAMapString)
			},
		},
		"correct grit_runner_managers Terraform State output": {
			input: `{"format_version": "1", "values": {"outputs": {"grit_runner_managers": {"value": {"runner-manager-1": {"instance_name": "instance-name", "address": "localhost:1234"}}}}}}`,
			assertMap: func(t *testing.T, m map[string]any) {
				if assert.Contains(t, m, "runner-manager-1") {
					rm, ok := (m["runner-manager-1"]).(map[string]any)
					require.True(t, ok)

					if assert.Contains(t, rm, "instance_name") {
						assert.Equal(t, rm["instance_name"], "instance-name")
					}

					if assert.Contains(t, rm, "address") {
						assert.Equal(t, rm["address"], "localhost:1234")
					}
				}
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			m, err := readRunnerManagersMap([]byte(tt.input))
			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, tt.assertMap, "assertMap must be defined in test definition")
			tt.assertMap(t, m)
		})
	}
}
