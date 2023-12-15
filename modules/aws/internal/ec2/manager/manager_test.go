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
				"runner_token":                           "runnerToken",
				"executor":                               "docker-autoscaler",
				"gitlab_url":                             "https://gitlab.com",
				"fleeting_service":                       "ec2",
				"fleeting_service_account_access_key_id": "access_key_id",
				"fleeting_service_account_secret_access_key": "secret_access_key",
				"ssh_key_pem":           "ssh_key_pem",
				"ssh_key_pem_name":      "ssh_key_pem_name",
				"aws_asg_name":          "aws_asg_name",
				"capacity_per_instance": 1,
				"scale_min":             1,
				"scale_max":             1,
				"vpc_id":                "1234",
				"subnet_id":             "1234",
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
