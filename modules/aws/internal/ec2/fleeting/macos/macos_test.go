package aws_ec2_fleeting_macos

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2FleetingMacos(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"aws_cloudformation_stack.jobs-cloudformation-stack",
		"aws_licensemanager_license_configuration.license-config",
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
				"name":              name,
			},
			expectedModules: expectedModules,
		},
		"override defaults": {
			moduleVars: map[string]interface{}{
				"aws_vpc_cidr":                   "10.0.0.0/24",
				"asg_ami_id":                     "ami-0fcd5ff1c92b00231",
				"asg_instance_type":              "mac2.metal",
				"asg_subnet_cidr":                "10.0.0.0/24",
				"name":                           name,
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
