package gorp

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools/tf"
)

func TestPerfectFit(t *testing.T) {
	modulesDir := filepath.Join("..", "modules")
	deps, err := tf.LoadModuleDependencies(modulesDir)
	require.NoError(t, err)
	for _, d := range deps {
		for _, o := range d.DependantModule.Outputs {
			assert.NotEmptyf(t, o.Type,
				"want output %v in module %v/%v to have a type. got none",
				o.Name,
				d.DependantModule.Provider,
				d.DependantModule.Name,
			)
		}
		assert.Equalf(t, d.DependantVariable.Type, d.DependencyOutputsType,
			"want variable %v in module %v/%v of type %v. got %v",
			d.DependantVariable.Name,
			d.DependantModule.Provider,
			d.DependantModule.Name,
			d.DependencyOutputsType,
			d.DependantVariable.Type,
		)
	}
}
