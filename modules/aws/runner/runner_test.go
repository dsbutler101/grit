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
					"subnet_ids": []string{"12345"},
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
				"s3_cache": map[string]interface{}{
					"enabled":           false,
					"server_address":    "s3.amazonaws.com",
					"bucket_name":       "bucket",
					"bucket_location":   "region",
					"access_key_id":     "access-key-id",
					"secret_access_key": "secret-access-key",
				},
				"service":               "ec2",
				"executor":              "docker-autoscaler",
				"scale_min":             2,
				"scale_max":             10,
				"idle_percentage":       10,
				"capacity_per_instance": 1,
				"security_group_ids":    []string{"123456"},
				"default_docker_image":  "ubuntu:latest",
				"associate_public_ip_address": true,
				"instance_type": "t2.micro",
				"encrypted": false,
				"kms_key_id": "",
				"volume_size": 8,
				"volume_type": "gp2",
				"throughput": 0,
				"usage_logger": map[string]interface{}{
					"enabled": true,
					"log_dir": "/var/log/usage",
					"custom_labels": map[string]interface{}{
						"stack_os": "os-name",
					},
				},
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
