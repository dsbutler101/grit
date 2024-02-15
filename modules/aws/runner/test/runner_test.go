package runner

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2Fleeting(t *testing.T) {
	name := test_tools.JobName(t)

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"ec2-docker-autoscaler runner_token provided": {
			moduleVars: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":   name,
					"labels": map[string]interface{}{},
				},
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
					"username":               "abcde",
				},
				"iam": map[string]interface{}{
					"fleeting_access_key_id":     "12345",
					"fleeting_secret_access_key": "abcde",
				},
				"service":            "ec2",
				"security_group_ids": []string{"123456"},
			},
			expectedModules: []string{
				"module.ec2[0].aws_instance.runner-manager",
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
