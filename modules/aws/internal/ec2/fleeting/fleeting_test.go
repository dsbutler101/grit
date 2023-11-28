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
		"module.common.aws_internet_gateway.internet-access",
		"module.common.aws_key_pair.jobs-key-pair",
		"module.common.aws_launch_template.fleeting-asg-template",
		"module.common.aws_route.internet-route",
		"module.common.aws_security_group.jobs-security-group",
		"module.common.aws_subnet.jobs-vpc-subnet",
		"module.common.aws_vpc.jobs-vpc",
		"module.common.data.aws_route_table.jobs-route-table",
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
		"linux fleet with all defaults": {
			moduleVars: map[string]interface{}{
				"fleeting_os":   "linux",
				"ami":           "ami-05012401516a40259",
				"instance_type": "t2.medium",
				"name":          name,
			},
			expectedModules: baseExpectedModules,
		},
		"linux fleet override defaults": {
			moduleVars: map[string]interface{}{
				"fleeting_os":               "linux",
				"ami":                       "ami-05012401516a40259",
				"instance_type":             "t2.medium",
				"name":                      name + "-linux-no-default",
				"aws_vpc_cidr":              "10.0.0.0/16",
				"gitlab_url":                "https://custom.gitlab.com/",
				"gitlab_runner_description": "eXtra GRIT",
				"gitlab_runner_tags":        []string{"grit-tag"},
				"scale_min":                 1,
				"scale_max":                 2,
				"idle_percentage":           10,
			},
			expectedModules: baseExpectedModules,
		},
		"mac fleet with all defaults": {
			moduleVars: map[string]interface{}{
				"fleeting_os":   "macos",
				"ami":           "ami-0fcd5ff1c92b00231",
				"instance_type": "mac2.metal",
				"name":          name + "-macos-default",
			},
			expectedModules: macExpectedModules,
		},
		"mac fleet override defaults": {
			moduleVars: map[string]interface{}{
				"fleeting_os":               "macos",
				"ami":                       "ami-0fcd5ff1c92b00231",
				"instance_type":             "mac2.metal",
				"name":                      name + "-macos-no-default",
				"aws_vpc_cidr":              "10.0.0.0/16",
				"gitlab_url":                "https://custom.gitlab.com/",
				"gitlab_runner_description": "eXtra GRIT",
				"gitlab_runner_tags":        []string{"grit-tag"},
				"scale_min":                 1,
				"scale_max":                 2,
				"idle_percentage":           10,
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
