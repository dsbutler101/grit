package runner

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestK8sRunner(t *testing.T) {
	plan := test_tools.Plan(t, map[string]any{
		"url":       "some-gitlab-url",
		"token":     "some-runner-registration-token",
		"namespace": "some-runner-namespace",
		"name":      "some-runner-name",
	})

	test_tools.AssertProviderConfigExists(t, plan, "kubectl")

	test_tools.AssertWithPlan(t, plan, []string{
		"kubectl_manifest.manifest",
		"kubectl_manifest.token_secret",
	})
}
