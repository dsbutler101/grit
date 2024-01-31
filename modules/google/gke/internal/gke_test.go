package manager

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestGKE(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"google_container_cluster.primary",
		"google_container_node_pool.primary_nodes",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"create gke": {
			moduleVars: map[string]interface{}{
				"name":              name,
				"labels":            map[string]string{"env": "another place"},
				"google_region":     "us-east1",
				"google_zone":       "us-east1-b",
				"nodes_count":       10,
				"node_machine_type": "e2-micro",
				"vpc": map[string]interface{}{
					"id":        "my-vpc",
					"subnet_id": "my-subnet",
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
