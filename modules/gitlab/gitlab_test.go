package gitlab

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/common"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2Fleeting(t *testing.T) {
	name := test_tools.JobName(t)

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"gitlab project runner": {
			moduleVars: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
				"url":                "https://gitlab.com",
				"runner_description": "my new GRIT runner",
				"project_id":         common.GritEndToEndTestProjectID,
				"runner_tags":        []string{"job-tag"},
			},
			expectedModules: []string{"gitlab_user_runner.primary"},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
