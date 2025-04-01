package manager

import (
	"maps"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2Manager(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"aws_instance.runner_manager",
	}

	variables := func(overrides map[string]any) map[string]any {
		required := map[string]any{
			"gitlab": map[string]any{
				"runner_token": "tokenString",
				"url":          "gitlab.com",
			},
			"vpc": map[string]any{
				"id":         "12345",
				"subnet_ids": []string{"12345"},
			},
			"fleeting": map[string]any{
				"autoscaling_group_name": "abcde",
				"ssh_key_pem_name":       "abcde",
				"ssh_key_pem":            "abcde",
				"username":               "abcde",
			},
			"iam": map[string]any{
				"fleeting_access_key_id":     "12345",
				"fleeting_secret_access_key": "abcde",
			},
			"s3_cache": map[string]any{
				"enabled":           false,
				"server_address":    "",
				"bucket_name":       "",
				"bucket_location":   "",
				"access_key_id":     "",
				"secret_access_key": "",
			},
			"executor":                    "docker-autoscaler",
			"capacity_per_instance":       1,
			"scale_min":                   1,
			"scale_max":                   1,
			"max_use_count":               10,
			"idle_percentage":             0.0,
			"idle_time":                   "20m0s",
			"name":                        name,
			"labels":                      map[string]string{"env": "another place"},
			"security_group_ids":          []string{"123456"},
			"install_cloudwatch_agent":    false,
			"region":                      "us-east-1",
			"runner_repository":           "gitlab-runner",
			"runner_version":              "16.11.1-1",
			"aws_plugin_version":          "0.5.0",
			"privileged":                  false,
			"cloudwatch_agent_json":       "ewogICJhZ2VudCI6IHsKICAgICJtZXRyaWNzX2NvbGxlY3Rpb25faW50ZXJ2YWwiOiA2MCwKICAgICJsb2dmaWxlIjogIi9vcHQvYXdzL2FtYXpvbi1jbG91ZHdhdGNoLWFnZW50L2xvZ3MvYW1hem9uLWNsb3Vkd2F0Y2gtYWdlbnQubG9nIiwKICAgICJkZWJ1ZyI6IGZhbHNlLAogICAgInJ1bl9hc191c2VyIjogImN3YWdlbnQiCiAgfSwKICAibG9ncyI6IHsKICAgICJsb2dzX2NvbGxlY3RlZCI6IHsKICAgICAgImZpbGVzIjogewogICAgICAgICJjb2xsZWN0X2xpc3QiOiBbCiAgICAgICAgICB7CiAgICAgICAgICAgICJmaWxlX3BhdGgiOiAiL3Zhci9sb2cvc3lzbG9nIiwKICAgICAgICAgICAgImxvZ19ncm91cF9uYW1lIjogIkZsZWV0aW5nLUxvZ3MiLAogICAgICAgICAgICAibG9nX3N0cmVhbV9uYW1lIjogIkZsZWV0aW5nLVN5c2xvZy1TdHJlYW0iLAogICAgICAgICAgICAidGltZXN0YW1wX2Zvcm1hdCI6ICIlSDogJU06ICVTJXklYiUtZCIKICAgICAgICAgIH0sCgkgIHsKICAgICAgICAgICAgImZpbGVfcGF0aCI6ICIvdmFyL2xvZy9jbG91ZC1pbml0LW91dHB1dC5sb2ciLAogICAgICAgICAgICAibG9nX2dyb3VwX25hbWUiOiAiRmxlZXRpbmctTG9ncyIsCiAgICAgICAgICAgICJsb2dfc3RyZWFtX25hbWUiOiAiRmxlZXRpbmctQ2xvdWRpbml0LVN0cmVhbSIsCiAgICAgICAgICAgICJ0aW1lc3RhbXBfZm9ybWF0IjogIiVIOiAlTTogJVMleSViJS1kIgogICAgICAgICAgfQoJXQogICAgICB9CiAgICB9CiAgfQp9Cg==",
			"instance_role_profile_name":  nil,
			"enable_metrics_export":       false,
			"metrics_export_endpoint":     "0.0.0.0:9402",
			"default_docker_image":        "ubuntu:latest",
			"associate_public_ip_address": true,
			"instance_type":               "t2.micro",
			"encrypted":                   false,
			"kms_key_id":                  "",
			"volume_size":                 8,
			"volume_type":                 "gp2",
			"throughput":                  0,
			"runner_manager_ami":          "ami-05012401516a40259",
		}
		maps.Copy(required, overrides)
		return required
	}

	testCases := map[string]struct {
		moduleVars      map[string]any
		expectedModules []string
		wantErr         bool
	}{
		"with required variables": {
			moduleVars:      variables(nil),
			expectedModules: expectedModules,
		},
		"with all optional variables": {
			moduleVars: variables(map[string]any{
				"usage_logger": map[string]any{
					"enabled": true,
					"log_dir": "/var/log/usage",
					"custom_labels": map[string]any{
						"stack_os": "os-name",
					},
				},
				"acceptable_durations": []map[string]any{{
					"periods":   []string{"1", "2", "3"},
					"threshold": "1",
					"timezone":  "UTC",
				}},
				"node_exporter": map[string]string{"enabled": "true", "port": "1234", "version": "1.2.3"},
				"runner_wrapper": map[string]any{
					"enabled":                     true,
					"process_termination_timeout": "1m",
					"socket_path":                 "tcp://foo",
				},
			}),
			expectedModules: expectedModules,
		},
		"vpc subnet_id and subnet_ids cannot be both specified": {
			moduleVars: variables(map[string]any{
				"vpc": map[string]any{
					"id":         "12345",
					"subnet_id":  "subnet-12345",
					"subnet_ids": []string{"subnet-12345"},
				},
			}),
			expectedModules: expectedModules,
			wantErr:         true,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			p, err := test_tools.PlanE(t, tc.moduleVars)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			test_tools.AssertWithPlan(t, p, expectedModules)
		})
	}
}
