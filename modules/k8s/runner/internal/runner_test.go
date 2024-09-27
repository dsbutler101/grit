package runner

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestK8sRunner(t *testing.T) {
	testCases := map[string]struct {
		vars            map[string]any
		expectedModules []string
		wantErr         bool
	}{
		"runner with config template": {
			vars: map[string]any{
				"url":              "some-gitlab-url",
				"token":            "some-runner-registration-token",
				"namespace":        "some-runner-namespace",
				"name":             "some-runner-name",
				"concurrent":       5,
				"check_interval":   10,
				"locked":           false,
				"protected":        false,
				"run_untagged":     true,
				"runner_tags":      []string{},
				"config_template":  "[[runners]]\\n		[runners.kubernetes]\\n		image=\\\"alpine\\\"",
				"runner_image":     "",
				"helper_image":     "",
				"pod_spec_patches": []map[string]string{},
			},
			expectedModules: []string{
				"kubectl_manifest.manifest",
				"kubectl_manifest.token_secret",
				"kubectl_manifest.config_template[0]",
			},
		},
		"runner with faulty config template": {
			vars: map[string]any{
				"url":              "some-gitlab-url",
				"token":            "some-runner-registration-token",
				"namespace":        "some-runner-namespace",
				"name":             "some-runner-name",
				"concurrent":       5,
				"check_interval":   10,
				"locked":           false,
				"protected":        false,
				"run_untagged":     true,
				"runner_tags":      []string{},
				"config_template":  "image=\\\"alpine\\\"",
				"runner_image":     "",
				"helper_image":     "",
				"pod_spec_patches": []map[string]string{},
			},
			wantErr: true,
		},
		"runner without config template": {
			vars: map[string]any{
				"url":              "some-gitlab-url",
				"token":            "some-runner-registration-token",
				"namespace":        "some-runner-namespace",
				"name":             "some-runner-name",
				"concurrent":       5,
				"check_interval":   10,
				"locked":           false,
				"protected":        false,
				"run_untagged":     true,
				"runner_tags":      []string{},
				"config_template":  "",
				"runner_image":     "",
				"helper_image":     "",
				"pod_spec_patches": []map[string]string{},
			},
			expectedModules: []string{
				"kubectl_manifest.manifest",
				"kubectl_manifest.token_secret",
			},
		},
		"runner with faulty pod spec": {
			vars: map[string]any{
				"url":             "some-gitlab-url",
				"token":           "some-runner-registration-token",
				"namespace":       "some-runner-namespace",
				"name":            "some-runner-name",
				"concurrent":      5,
				"check_interval":  10,
				"locked":          false,
				"protected":       false,
				"run_untagged":    true,
				"runner_tags":     []string{},
				"config_template": "",
				"runner_image":    "",
				"helper_image":    "",
				"pod_spec_patches": []map[string]string{
					{
						"name":      "custom-deadline-seconds",
						"patchType": "merge",
					},
				},
			},
			wantErr: true,
		},
		"runner with active deadline seconds pod spec": {
			vars: map[string]any{
				"url":             "some-gitlab-url",
				"token":           "some-runner-registration-token",
				"namespace":       "some-runner-namespace",
				"name":            "some-runner-name",
				"concurrent":      5,
				"check_interval":  10,
				"locked":          false,
				"protected":       false,
				"run_untagged":    true,
				"runner_tags":     []string{},
				"config_template": "",
				"runner_image":    "",
				"helper_image":    "",
				"pod_spec_patches": []map[string]string{
					{
						"name":      "custom-deadline-seconds",
						"patch":     "activeDeadlineSeconds: 300",
						"patchType": "merge",
					},
				},
			},
			expectedModules: []string{
				"kubectl_manifest.manifest",
				"kubectl_manifest.token_secret",
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			if tc.wantErr {
				test_tools.PlanAndAssertError(t, tc.vars, tc.wantErr)
				return
			}

			plan := test_tools.Plan(t, tc.vars)
			test_tools.AssertProviderConfigExists(t, plan, "kubectl")
			test_tools.AssertWithPlan(t, plan, tc.expectedModules)
		})
	}
}
