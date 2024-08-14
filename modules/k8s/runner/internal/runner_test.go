package runner

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestK8sRunnerWithConfigTemplate(t *testing.T) {
	plan := test_tools.Plan(t, map[string]any{
		"url":       "some-gitlab-url",
		"token":     "some-runner-registration-token",
		"namespace": "some-runner-namespace",
		"name":      "some-runner-name",
		"config_template": `shutdown_timeout = 100
[[runners]]
    [runners.kubernetes]
    image = "alpine"`,
	})

	test_tools.AssertProviderConfigExists(t, plan, "kubectl")

	test_tools.AssertWithPlan(t, plan, []string{
		"kubectl_manifest.manifest",
		"kubectl_manifest.token_secret",
		"kubectl_manifest.config_template[0]",
	})
}

func TestK8sRunnerWithoutConfigTemplate(t *testing.T) {
	plan := test_tools.Plan(t, map[string]any{
		"url":             "some-gitlab-url",
		"token":           "some-runner-registration-token",
		"namespace":       "some-runner-namespace",
		"name":            "some-runner-name",
		"config_template": "",
	})

	test_tools.AssertProviderConfigExists(t, plan, "kubectl")

	test_tools.AssertWithPlan(t, plan, []string{
		"kubectl_manifest.manifest",
		"kubectl_manifest.token_secret",
	})
}
