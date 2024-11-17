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
				"config_template": "[[runners]]\\n		[runners.kubernetes]\\n		image=\\\"alpine\\\"",
			},
			expectedModules: []string{
				"kubectl_manifest.manifest",
				"kubectl_manifest.token_secret",
				"kubectl_manifest.config_template[0]",
				"terraform_data.config_template",
				"terraform_data.token_secret",
				"terraform_data.envvars",
			},
		},
		"runner with faulty config template": {
			vars: map[string]any{
				"config_template": "image=\\\"alpine\\\"",
			},
			wantErr: true,
		},
		"runner without config template": {
			vars: map[string]any{
				"config_template": "",
			},
			expectedModules: []string{
				"kubectl_manifest.manifest",
				"kubectl_manifest.token_secret",
				"terraform_data.config_template",
				"terraform_data.token_secret",
				"terraform_data.envvars",
			},
		},
		"runner with faulty pod spec": {
			vars: map[string]any{
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
				"terraform_data.config_template",
				"terraform_data.token_secret",
				"terraform_data.envvars",
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			defaultVars := map[string]any{
				"metadata": map[string]any{
					"name":        "some-runner-name",
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
				"namespace": "some-runner-namespace",
				"gitlab": map[string]any{
					"url":          "some-gitlab-url",
					"runner_token": "some-runner-registration-token",
				},
			}

			for k, v := range defaultVars {
				tc.vars[k] = v
			}

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
