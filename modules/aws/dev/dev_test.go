package aws_dev

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSDev(t *testing.T) {
	name := test_tools.JobName(t)

	macExpectedModules := []string{
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_autoscaling_group.fleeting-asg",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_iam_access_key.fleeting-service-account-key",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_iam_policy.fleeting-service-account-policy",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_iam_user.fleeting-service-account",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_iam_user_policy_attachment.fleeting-service-account-attach",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_internet_gateway.internet-access",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_key_pair.jobs-key-pair",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_launch_template.fleeting-asg-template",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_route.internet-route",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_security_group.jobs-security-group",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_subnet.jobs-vpc-subnet",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.aws_vpc.jobs-vpc",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.data.aws_route_table.jobs-route-table",
		"module.ec2[0].module.macos[0].module.instance_group.module.common.tls_private_key.aws-jobs-private-key",
		"module.ec2[0].module.macos[0].module.instance_group.module.macos[0].aws_cloudformation_stack.jobs-cloudformation-stack",
		"module.ec2[0].module.macos[0].module.instance_group.module.macos[0].aws_licensemanager_license_configuration.license-config",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"ec2 linux fleet": {
			moduleVars: map[string]interface{}{
				"fleeting_service": "ec2",
				"fleeting_os":      "linux",
				"name":             name,
			},
			// linux currently not implemented in dev, should return empty expectedModules
			expectedModules: []string{},
		},
		"ec2 windows fleet": {
			moduleVars: map[string]interface{}{
				"fleeting_service": "ec2",
				"fleeting_os":      "linux",
				"name":             name,
			},
			// windows currently not implemented in dev, should return empty expectedModules
			expectedModules: []string{},
		},
		"ec2 macos fleet": {
			moduleVars: map[string]interface{}{
				"fleeting_service": "ec2",
				"fleeting_os":      "macos",
				"name":             name,
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
