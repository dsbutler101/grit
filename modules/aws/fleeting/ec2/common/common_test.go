package aws_internal_ec2_fleeting_common

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2FleetingCommon(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"aws_autoscaling_group.fleeting_asg",
		"aws_key_pair.jobs_key_pair",
		"aws_launch_template.fleeting_asg_template",
		"tls_private_key.aws_jobs_private_key",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"common fleet": {
			moduleVars: map[string]interface{}{
				"license_arn":                      "",
				"jobs_host_resource_group_outputs": map[string]string{},
				"scale_min":                        0,
				"scale_max":                        10,
				"storage_size":                     500,
				"storage_type":                     "gp3",
				"storage_throughput":               125,
				"ephemeral_runner_ami":             "ami-0fcd5ff1c92b00231",
				"instance_type":                    "mac2.metal",
				"labels":                           map[string]string{"env": "another-place"},
				"subnet_ids":                       []string{"12345"},
				"name":                             name + "-macos",
				"security_group_ids":               []string{"123456"},
				"install_cloudwatch_agent":         false,
				"cloudwatch_agent_json":            "ewogICJhZ2VudCI6IHsKICAgICJtZXRyaWNzX2NvbGxlY3Rpb25faW50ZXJ2YWwiOiA2MCwKICAgICJsb2dmaWxlIjogIi9vcHQvYXdzL2FtYXpvbi1jbG91ZHdhdGNoLWFnZW50L2xvZ3MvYW1hem9uLWNsb3Vkd2F0Y2gtYWdlbnQubG9nIiwKICAgICJkZWJ1ZyI6IGZhbHNlLAogICAgInJ1bl9hc191c2VyIjogImN3YWdlbnQiCiAgfSwKICAibG9ncyI6IHsKICAgICJsb2dzX2NvbGxlY3RlZCI6IHsKICAgICAgImZpbGVzIjogewogICAgICAgICJjb2xsZWN0X2xpc3QiOiBbCiAgICAgICAgICB7CiAgICAgICAgICAgICJmaWxlX3BhdGgiOiAiL3Zhci9sb2cvc3lzbG9nIiwKICAgICAgICAgICAgImxvZ19ncm91cF9uYW1lIjogIkZsZWV0aW5nLUxvZ3MiLAogICAgICAgICAgICAibG9nX3N0cmVhbV9uYW1lIjogIkZsZWV0aW5nLVN5c2xvZy1TdHJlYW0iLAogICAgICAgICAgICAidGltZXN0YW1wX2Zvcm1hdCI6ICIlSDogJU06ICVTJXklYiUtZCIKICAgICAgICAgIH0sCgkgIHsKICAgICAgICAgICAgImZpbGVfcGF0aCI6ICIvdmFyL2xvZy9jbG91ZC1pbml0LW91dHB1dC5sb2ciLAogICAgICAgICAgICAibG9nX2dyb3VwX25hbWUiOiAiRmxlZXRpbmctTG9ncyIsCiAgICAgICAgICAgICJsb2dfc3RyZWFtX25hbWUiOiAiRmxlZXRpbmctQ2xvdWRpbml0LVN0cmVhbSIsCiAgICAgICAgICAgICJ0aW1lc3RhbXBfZm9ybWF0IjogIiVIOiAlTTogJVMleSViJS1kIgogICAgICAgICAgfQoJXQogICAgICB9CiAgICB9CiAgfQp9Cg==",
				"instance_role_profile_name":       nil,
				"ebs_encryption":                   false,
				"kms_key_arn":                      "",
				"mixed_instances_policy":           nil,
				"node_exporter":                    map[string]string{"enabled": "true", "port": "1234", "version": "1.2.3"},
			},
			expectedModules: expectedModules,
		},
		"common fleet - node exporter disabled": {
			moduleVars: map[string]interface{}{
				"license_arn":                      "",
				"jobs_host_resource_group_outputs": map[string]string{},
				"scale_min":                        0,
				"scale_max":                        10,
				"storage_size":                     500,
				"storage_type":                     "gp3",
				"storage_throughput":               125,
				"ephemeral_runner_ami":             "ami-0fcd5ff1c92b00231",
				"instance_type":                    "mac2.metal",
				"labels":                           map[string]string{"env": "another-place"},
				"subnet_ids":                       []string{"12345"},
				"name":                             name + "-macos",
				"security_group_ids":               []string{"123456"},
				"install_cloudwatch_agent":         false,
				"cloudwatch_agent_json":            "ewogICJhZ2VudCI6IHsKICAgICJtZXRyaWNzX2NvbGxlY3Rpb25faW50ZXJ2YWwiOiA2MCwKICAgICJsb2dmaWxlIjogIi9vcHQvYXdzL2FtYXpvbi1jbG91ZHdhdGNoLWFnZW50L2xvZ3MvYW1hem9uLWNsb3Vkd2F0Y2gtYWdlbnQubG9nIiwKICAgICJkZWJ1ZyI6IGZhbHNlLAogICAgInJ1bl9hc191c2VyIjogImN3YWdlbnQiCiAgfSwKICAibG9ncyI6IHsKICAgICJsb2dzX2NvbGxlY3RlZCI6IHsKICAgICAgImZpbGVzIjogewogICAgICAgICJjb2xsZWN0X2xpc3QiOiBbCiAgICAgICAgICB7CiAgICAgICAgICAgICJmaWxlX3BhdGgiOiAiL3Zhci9sb2cvc3lzbG9nIiwKICAgICAgICAgICAgImxvZ19ncm91cF9uYW1lIjogIkZsZWV0aW5nLUxvZ3MiLAogICAgICAgICAgICAibG9nX3N0cmVhbV9uYW1lIjogIkZsZWV0aW5nLVN5c2xvZy1TdHJlYW0iLAogICAgICAgICAgICAidGltZXN0YW1wX2Zvcm1hdCI6ICIlSDogJU06ICVTJXklYiUtZCIKICAgICAgICAgIH0sCgkgIHsKICAgICAgICAgICAgImZpbGVfcGF0aCI6ICIvdmFyL2xvZy9jbG91ZC1pbml0LW91dHB1dC5sb2ciLAogICAgICAgICAgICAibG9nX2dyb3VwX25hbWUiOiAiRmxlZXRpbmctTG9ncyIsCiAgICAgICAgICAgICJsb2dfc3RyZWFtX25hbWUiOiAiRmxlZXRpbmctQ2xvdWRpbml0LVN0cmVhbSIsCiAgICAgICAgICAgICJ0aW1lc3RhbXBfZm9ybWF0IjogIiVIOiAlTTogJVMleSViJS1kIgogICAgICAgICAgfQoJXQogICAgICB9CiAgICB9CiAgfQp9Cg==",
				"instance_role_profile_name":       nil,
				"ebs_encryption":                   false,
				"kms_key_arn":                      "",
				"mixed_instances_policy":           nil,
				"node_exporter":                    map[string]string{"enabled": "false"},
			},
			expectedModules: expectedModules,
		},
		"common fleet mix instance policy": {
			moduleVars: map[string]interface{}{
				"license_arn":                      "",
				"jobs_host_resource_group_outputs": map[string]string{},
				"scale_min":                        0,
				"scale_max":                        10,
				"storage_size":                     500,
				"storage_type":                     "gp3",
				"storage_throughput":               125,
				"ephemeral_runner_ami":             "ami-0fcd5ff1c92b00231",
				"instance_type":                    "mac2.metal",
				"labels":                           map[string]string{"env": "another-place"},
				"subnet_ids":                       []string{"12345"},
				"name":                             name + "-macos",
				"security_group_ids":               []string{"123456"},
				"install_cloudwatch_agent":         false,
				"cloudwatch_agent_json":            "ewogICJhZ2VudCI6IHsKICAgICJtZXRyaWNzX2NvbGxlY3Rpb25faW50ZXJ2YWwiOiA2MCwKICAgICJsb2dmaWxlIjogIi9vcHQvYXdzL2FtYXpvbi1jbG91ZHdhdGNoLWFnZW50L2xvZ3MvYW1hem9uLWNsb3Vkd2F0Y2gtYWdlbnQubG9nIiwKICAgICJkZWJ1ZyI6IGZhbHNlLAogICAgInJ1bl9hc191c2VyIjogImN3YWdlbnQiCiAgfSwKICAibG9ncyI6IHsKICAgICJsb2dzX2NvbGxlY3RlZCI6IHsKICAgICAgImZpbGVzIjogewogICAgICAgICJjb2xsZWN0X2xpc3QiOiBbCiAgICAgICAgICB7CiAgICAgICAgICAgICJmaWxlX3BhdGgiOiAiL3Zhci9sb2cvc3lzbG9nIiwKICAgICAgICAgICAgImxvZ19ncm91cF9uYW1lIjogIkZsZWV0aW5nLUxvZ3MiLAogICAgICAgICAgICAibG9nX3N0cmVhbV9uYW1lIjogIkZsZWV0aW5nLVN5c2xvZy1TdHJlYW0iLAogICAgICAgICAgICAidGltZXN0YW1wX2Zvcm1hdCI6ICIlSDogJU06ICVTJXklYiUtZCIKICAgICAgICAgIH0sCgkgIHsKICAgICAgICAgICAgImZpbGVfcGF0aCI6ICIvdmFyL2xvZy9jbG91ZC1pbml0LW91dHB1dC5sb2ciLAogICAgICAgICAgICAibG9nX2dyb3VwX25hbWUiOiAiRmxlZXRpbmctTG9ncyIsCiAgICAgICAgICAgICJsb2dfc3RyZWFtX25hbWUiOiAiRmxlZXRpbmctQ2xvdWRpbml0LVN0cmVhbSIsCiAgICAgICAgICAgICJ0aW1lc3RhbXBfZm9ybWF0IjogIiVIOiAlTTogJVMleSViJS1kIgogICAgICAgICAgfQoJXQogICAgICB9CiAgICB9CiAgfQp9Cg==",
				"instance_role_profile_name":       nil,
				"ebs_encryption":                   false,
				"kms_key_arn":                      "",
				"mixed_instances_policy": map[string]interface{}{
					"override": []map[string]string{
						{"instance_type": "t3.nano"},
						{"instance_type": "t3.medium"},
					},
				},
				"node_exporter": map[string]string{"enabled": "true", "port": "1234", "version": "1.2.3"},
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
