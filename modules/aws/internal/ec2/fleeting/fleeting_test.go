package aws_ec2_fleeting

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2Fleeting(t *testing.T) {
	name := test_tools.JobName(t)

	baseExpectedModules := []string{
		"module.common.aws_autoscaling_group.fleeting-asg",
		"module.common.aws_iam_access_key.fleeting-service-account-key",
		"module.common.aws_iam_policy.fleeting-service-account-policy",
		"module.common.aws_iam_user.fleeting-service-account",
		"module.common.aws_iam_user_policy_attachment.fleeting-service-account-attach",
		"module.common.aws_key_pair.jobs-key-pair",
		"module.common.aws_launch_template.fleeting-asg-template",
		"module.common.aws_security_group.jobs-security-group",
		"module.common.tls_private_key.aws-jobs-private-key",
	}

	macExpectedModules := append(baseExpectedModules, []string{
		"module.macos[0].aws_cloudformation_stack.jobs-cloudformation-stack",
		"module.macos[0].aws_licensemanager_license_configuration.license-config",
	}...)

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"linux fleet": {
			moduleVars: map[string]interface{}{
				"fleeting_os":                          "linux",
				"ami":                                  "ami-05012401516a40259",
				"scale_min":                            1,
				"scale_max":                            2,
				"idle_percentage":                      10,
				"asg_storage_size":                     500,
				"asg_storage_type":                     "gp3",
				"asg_storage_throughput":               125,
				"instance_type":                        "t2.medium",
				"vpc_id":                               "1234",
				"subnet_id":                            "1234",
				"macos_required_license_count_per_asg": 0,
				"macos_cores_per_license":              0,
				"labels":                               map[string]string{"env": "another place"},
				"name":                                 name + "-linux-no-default",
			},
			expectedModules: baseExpectedModules,
		},
		"mac fleet": {
			moduleVars: map[string]interface{}{
				"fleeting_os":                          "macos",
				"ami":                                  "ami-0fcd5ff1c92b00231",
				"instance_type":                        "mac2.metal",
				"scale_min":                            1,
				"scale_max":                            2,
				"idle_percentage":                      10,
				"asg_storage_size":                     500,
				"asg_storage_type":                     "gp3",
				"asg_storage_throughput":               125,
				"vpc_id":                               "1234",
				"subnet_id":                            "1234",
				"macos_required_license_count_per_asg": 10,
				"macos_cores_per_license":              4,
				"labels":                               map[string]string{"env": "another place"},
				"name":                                 name + "-macos-no-default",
			},
			expectedModules: macExpectedModules,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
