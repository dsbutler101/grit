package examples

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestExamples(t *testing.T) {
	entries, err := os.ReadDir("./")
	require.NoError(t, err)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		t.Run(e.Name(), func(t *testing.T) {
			require.NoError(t, err)
			options := &terraform.Options{
				TerraformBinary: "terraform",
				TerraformDir:    e.Name(),
			}
			_, err := terraform.InitAndPlanE(t, options)
			require.NoError(t, err)
		})
	}
}
