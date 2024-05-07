package manager

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2Manager(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"aws_instance.runner-manager",
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
					"username":               "abcde",
				},
				"iam": map[string]interface{}{
					"fleeting_access_key_id":     "12345",
					"fleeting_secret_access_key": "abcde",
				},
				"executor":                   "docker-autoscaler",
				"capacity_per_instance":      1,
				"scale_min":                  1,
				"scale_max":                  1,
				"name":                       name,
				"labels":                     map[string]string{"env": "another place"},
				"security_group_ids":         []string{"123456"},
				"install_cloudwatch_agent":   false,
				"region":                     "us-east-1",
				"runner_version":             "16.11.1-1",
				"aws_plugin_version":         "0.5.0",
				"privileged":                 false,
				"cloudwatch_agent_json":      "ewogICJhZ2VudCI6IHsKICAgICJtZXRyaWNzX2NvbGxlY3Rpb25faW50ZXJ2YWwiOiA2MCwKICAgICJsb2dmaWxlIjogIi9vcHQvYXdzL2FtYXpvbi1jbG91ZHdhdGNoLWFnZW50L2xvZ3MvYW1hem9uLWNsb3Vkd2F0Y2gtYWdlbnQubG9nIiwKICAgICJkZWJ1ZyI6IGZhbHNlLAogICAgInJ1bl9hc191c2VyIjogImN3YWdlbnQiCiAgfSwKICAibG9ncyI6IHsKICAgICJsb2dzX2NvbGxlY3RlZCI6IHsKICAgICAgImZpbGVzIjogewogICAgICAgICJjb2xsZWN0X2xpc3QiOiBbCiAgICAgICAgICB7CiAgICAgICAgICAgICJmaWxlX3BhdGgiOiAiL3Zhci9sb2cvc3lzbG9nIiwKICAgICAgICAgICAgImxvZ19ncm91cF9uYW1lIjogIkZsZWV0aW5nLUxvZ3MiLAogICAgICAgICAgICAibG9nX3N0cmVhbV9uYW1lIjogIkZsZWV0aW5nLVN5c2xvZy1TdHJlYW0iLAogICAgICAgICAgICAidGltZXN0YW1wX2Zvcm1hdCI6ICIlSDogJU06ICVTJXklYiUtZCIKICAgICAgICAgIH0sCgkgIHsKICAgICAgICAgICAgImZpbGVfcGF0aCI6ICIvdmFyL2xvZy9jbG91ZC1pbml0LW91dHB1dC5sb2ciLAogICAgICAgICAgICAibG9nX2dyb3VwX25hbWUiOiAiRmxlZXRpbmctTG9ncyIsCiAgICAgICAgICAgICJsb2dfc3RyZWFtX25hbWUiOiAiRmxlZXRpbmctQ2xvdWRpbml0LVN0cmVhbSIsCiAgICAgICAgICAgICJ0aW1lc3RhbXBfZm9ybWF0IjogIiVIOiAlTTogJVMleSViJS1kIgogICAgICAgICAgfQoJXQogICAgICB9CiAgICB9CiAgfQp9Cg==",
				"instance_role_profile_name": nil,
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
