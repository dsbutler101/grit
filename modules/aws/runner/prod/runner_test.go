package runner

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestRunner(t *testing.T) {
	name := test_tools.JobName(t)

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"ec2-docker-autoscaler runner_token provided and preexisting vpc": {
			moduleVars: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
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
				"service":               "ec2",
				"executor":              "docker-autoscaler",
				"scale_min":             2,
				"scale_max":             10,
				"idle_percentage":       10,
				"capacity_per_instance": 1,
				"security_group_ids":    []string{"123456"},
				"default_docker_image":  "ubuntu:latest",
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
