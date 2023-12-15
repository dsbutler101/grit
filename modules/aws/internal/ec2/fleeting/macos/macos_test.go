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
		"override defaults": {
			moduleVars: map[string]interface{}{
				"required_license_count_per_asg": 10,
				"cores_per_license":              4,
				"labels":                         map[string]string{"env": "another place"},
				"name":                           name,
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
