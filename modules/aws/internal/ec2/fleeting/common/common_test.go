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
		"aws_internet_gateway.internet-access",
		"aws_key_pair.jobs-key-pair",
		"aws_launch_template.fleeting-asg-template",
		"aws_route.internet-route",
		"aws_security_group.jobs-security-group",
		"aws_subnet.jobs-vpc-subnet",
		"aws_vpc.jobs-vpc",
		"data.aws_route_table.jobs-route-table",
		"tls_private_key.aws-jobs-private-key",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"do not override defaults": {
			moduleVars: map[string]interface{}{
				"aws_vpc_cidr":      "10.0.0.0/24",
				"asg_ami_id":        "ami-0fcd5ff1c92b00231",
				"asg_instance_type": "mac2.metal",
				"asg_subnet_cidr":   "10.0.0.0/24",
				"name":              name + "-common-default",
			},
			expectedModules: expectedModules,
		},
		"override defaults": {
			moduleVars: map[string]interface{}{
				"aws_vpc_cidr":                   "10.0.0.0/24",
				"asg_ami_id":                     "ami-0fcd5ff1c92b00231",
				"asg_instance_type":              "mac2.metal",
				"asg_subnet_cidr":                "10.0.0.0/24",
				"name":                           name + "-macos",
				"required_license_count_per_asg": 10,
				"cores_per_license":              4,
				"labels":                         map[string]string{"env": "another place"},
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
