package manager

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2Manager(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"aws_instance.runner-manager",
		"aws_security_group.manager_sg",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"manager override defaults": {
			moduleVars: map[string]interface{}{
				"gitlab": map[string]interface{}{
					"runner_token": "tokenString",
					"url":          "gitlab.com",
				},
				"vpc": map[string]interface{}{
					"id":        "12345",
					"subnet_id": "12345",
				},
				"fleeting": map[string]interface{}{
					"autoscaling_group_name": "abcde",
					"ssh_key_pem_name":       "abcde",
					"ssh_key_pem":            "abcde",
				},
				"iam": map[string]interface{}{
					"fleeting_access_key_id":     "12345",
					"fleeting_secret_access_key": "abcde",
				},
				"executor":              "docker-autoscaler",
				"capacity_per_instance": 1,
				"scale_min":             1,
				"scale_max":             1,
				"name":                  name,
			},
			expectedModules: expectedModules,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
