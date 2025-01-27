package vpc

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestVPC(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"google_compute_network.default",
		"google_compute_firewall.runner-manager-ingress-default",
		"google_compute_firewall.runner-manager-egress-default",
		`google_compute_subnetwork.subnetwork["subnet-1"]`,
		`google_compute_subnetwork.subnetwork["subnet-2"]`,
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"create vpc": {
			moduleVars: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{"env": "another-place"},
					"min_support": "experimental",
				},
				"google_region": "us-east1",
				"subnetworks": map[string]string{
					"subnet-1": "10.0.0.0/24",
					"subnet-2": "10.0.1.0/24",
				},
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
