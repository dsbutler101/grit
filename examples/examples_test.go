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

				// We inject some config via `EnvVars` rather than via `Vars`;
				// terraform will complain about variables being set on the commandline
				// when there is no such terraform variable in the root module (e.g. no
				// AWS module would have a `google_project` variable).
				EnvVars: map[string]string{
					"TF_VAR_google_project":    mustFromEmv(t, "GOOGLE_PROJECT"),
					"TF_VAR_google_region":     mustFromEmv(t, "GOOGLE_REGION"),
					"TF_VAR_google_zone":       mustFromEmv(t, "GOOGLE_ZONE"),
					"TF_VAR_gitlab_pat":        mustFromEmv(t, "GITLAB_TOKEN"),
					"TF_VAR_gitlab_project_id": "39258790",
				},
			}
			_, err := terraform.InitAndPlanE(t, options)
			require.NoError(t, err)
		})
	}
}

func mustFromEmv(t *testing.T, name string) string {
	t.Helper()

	value, ok := os.LookupEnv(name)
	require.True(t, ok, "expected environment variable %q to be set, but is not", name)

	return value
}
