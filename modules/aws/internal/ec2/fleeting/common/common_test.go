package aws_internal_ec2_fleeting_common

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2FleetingCommon(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"aws_autoscaling_group.fleeting-asg",
		"aws_iam_access_key.fleeting-service-account-key",
		"aws_iam_policy.fleeting-service-account-policy",
		"aws_iam_user.fleeting-service-account",
		"aws_iam_user_policy_attachment.fleeting-service-account-attach",
		"aws_key_pair.jobs-key-pair",
		"aws_launch_template.fleeting-asg-template",
		"aws_security_group.jobs-security-group",
		"tls_private_key.aws-jobs-private-key",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"common fleet": {
			moduleVars: map[string]interface{}{
				"license_arn":                      "",
				"jobs-host-resource-group-outputs": map[string]string{},
				"scale_min":                        0,
				"scale_max":                        10,
				"idle_percentage":                  10,
				"asg_storage_size":                 500,
				"asg_storage_type":                 "gp3",
				"asg_storage_throughput":           125,
				"asg_ami_id":                       "ami-0fcd5ff1c92b00231",
				"asg_instance_type":                "mac2.metal",
				"labels":                           map[string]string{"env": "another place"},
				"vpc_id":                           "1234",
				"subnet_id":                        "12345",
				"name":                             name + "-macos",
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
