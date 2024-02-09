package securitygroups

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSSecurityGroups(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"aws_security_group.manager_sg",
		"aws_security_group.jobs_security_group",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"common fleet": {
			moduleVars: map[string]interface{}{
				"labels": map[string]string{"env": "another place"},
				"vpc_id": "1234",
				"name":   name,
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
