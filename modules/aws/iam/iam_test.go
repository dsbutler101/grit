package runner

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestIam(t *testing.T) {
	name := test_tools.JobName(t)

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"iam": {
			moduleVars: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
			},
			expectedModules: []string{
				"aws_iam_access_key.fleeting-service-account-key",
				"aws_iam_policy.fleeting-service-account-policy",
				"aws_iam_user.fleeting-service-account",
				"aws_iam_user_policy_attachment.fleeting-service-account-attach",
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
