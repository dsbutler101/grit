package aws_dev

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSDev(t *testing.T) {
	name := test_tools.JobName(t)

	baseExpectedModules := []string{
		"module.dev-module.module.ec2-instance-group[0].module.common.aws_autoscaling_group.fleeting-asg",
		"module.dev-module.module.ec2-instance-group[0].module.common.aws_iam_access_key.fleeting-service-account-key",
		"module.dev-module.module.ec2-instance-group[0].module.common.aws_iam_policy.fleeting-service-account-policy",
		"module.dev-module.module.ec2-instance-group[0].module.common.aws_iam_user.fleeting-service-account",
		"module.dev-module.module.ec2-instance-group[0].module.common.aws_iam_user_policy_attachment.fleeting-service-account-attach",
		"module.dev-module.module.ec2-instance-group[0].module.common.aws_key_pair.jobs-key-pair",
		"module.dev-module.module.ec2-instance-group[0].module.common.aws_launch_template.fleeting-asg-template",
		"module.dev-module.module.ec2-instance-group[0].module.common.aws_security_group.jobs-security-group",
		"module.dev-module.module.ec2-instance-group[0].module.common.tls_private_key.aws-jobs-private-key",
		"module.dev-module.module.vpc[0].aws_internet_gateway.igw",
		"module.dev-module.module.vpc[0].aws_route.internet-route",
		"module.dev-module.module.vpc[0].aws_route_table.rtb_public",
		"module.dev-module.module.vpc[0].aws_route_table_association.rta_subnet_public",
		"module.dev-module.module.vpc[0].aws_subnet.jobs-vpc-subnet",
		"module.dev-module.module.vpc[0].aws_vpc.vpc",
	}

	vpcExpectedModules := []string{
		"module.dev-module.module.vpc[0].aws_internet_gateway.igw",
		"module.dev-module.module.vpc[0].aws_route.internet-route",
		"module.dev-module.module.vpc[0].aws_route_table.rtb_public",
		"module.dev-module.module.vpc[0].aws_route_table_association.rta_subnet_public",
		"module.dev-module.module.vpc[0].aws_subnet.jobs-vpc-subnet",
		"module.dev-module.module.vpc[0].aws_vpc.vpc",
	}

	macExpectedModules := append(baseExpectedModules, []string{
		"module.dev-module.module.ec2-instance-group[0].module.macos[0].aws_cloudformation_stack.jobs-cloudformation-stack",
		"module.dev-module.module.ec2-instance-group[0].module.macos[0].aws_licensemanager_license_configuration.license-config",
	}...)

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"ec2 linux fleet": {
			moduleVars: map[string]interface{}{
				"fleeting_service": "ec2",
				"fleeting_os":      "linux",
				"ami":              "ami-05012401516a40259",
				"instance_type":    "t2.medium",
				"name":             name,
			},
			// linux currently not implemented in dev, should return empty expectedModules
			expectedModules: baseExpectedModules,
		},
		"ec2 windows fleet": {
			moduleVars: map[string]interface{}{
				"fleeting_service": "ec2",
				"fleeting_os":      "windows",
				"ami":              "none",
				"instance_type":    "none",
				"name":             name,
			},
			// windows currently not implemented, only vpc will get created
			expectedModules: vpcExpectedModules,
		},
		"ec2 macos fleet": {
			moduleVars: map[string]interface{}{
				"fleeting_service": "ec2",
				"fleeting_os":      "macos",
				"ami":              "ami-0fcd5ff1c92b00231",
				"instance_type":    "mac2.metal",
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
