package manager

import (
	"maps"
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

type moduleVars = map[string]any

func defaultModuleVars(t *testing.T) moduleVars {
	return moduleVars{
		"name":        test_tools.JobName(t),
		"google_zone": "us-east1-b",
		"metadata": map[string]any{
			"labels": map[string]string{
				"some": "label",
			},
		},
		"deletion_protection": "false",
		"node_pools": map[string]any{
			"default": map[string]any{
				"node_count": 1,
				"node_config": map[string]any{
					"machine_type": "e2-micro",
					"taints": []map[string]string{
						{
							"key":    "key",
							"value":  "value",
							"effect": "NO_EXECUTE",
						},
					},
				},
			},
			"vpc": map[string]string{
				"id":        "",
				"subnet_id": "",
			},
			"autoscaling": map[string]any{
				"enabled":                     false,
				"autoscaling_profile":         "",
				"auto_provisioning_locations": []string{},
				"resource_limits": []map[string]any{
					{
						"resource_type": "cpu",
						"minimum":       1,
						"maximum":       10,
					},
					{
						"resource_type": "memory",
						"minimum":       1,
						"maximum":       10,
					},
				},
			},
		},
	}
}

func TestGKE(t *testing.T) {
	expectedModules := []string{
		"google_container_cluster.primary",
		`google_container_node_pool.linux_node_pool["default"]`,
		`google_container_node_pool.linux_node_pool["autoscaling"]`,
		`google_container_node_pool.linux_node_pool["vpc"]`,
	}

	testCases := map[string]struct {
		moduleVars      moduleVars
		expectedModules []string
	}{
		"create-gke-happy": {
			moduleVars: moduleVars{
				"metadata": map[string]any{
					"labels":      map[string]string{"env": "another-place"},
					"min_support": "experimental",
					"name":        "gke",
				},
				"vpc": map[string]any{
					"id":        "my-vpc",
					"subnet_id": "my-subnet",
				},
			},
			expectedModules: expectedModules,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			mvs := mergeModuleVars(
				defaultModuleVars(t),
				tc.moduleVars,
			)
			test_tools.PlanAndAssert(t, mvs, tc.expectedModules)
		})
	}
}

func TestGKEPlanErrors(t *testing.T) {
	testCases := map[string]struct {
		moduleVars     map[string]interface{}
		shouldNotError bool
	}{
		"invalid-label-value": {
			moduleVars: map[string]interface{}{
				"metadata": map[string]any{
					"labels": map[string]string{
						"label1": "google does not like this because it has spaces",
						"label2": "this-is-fine",
					},
					"min_support": "experimental",
					"name":        "gke",
				},
			},
		},
		"valid-label-value": {
			moduleVars: moduleVars{
				"metadata": map[string]any{
					"labels": map[string]string{
						"label1": "this-is-fine",
						"label2": "this.is_fine-too",
					},
					"min_support": "experimental",
					"name":        "gke",
				},
			},
			shouldNotError: true,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			mvs := mergeModuleVars(
				defaultModuleVars(t),
				tc.moduleVars,
			)
			test_tools.PlanAndAssertError(t, mvs, !tc.shouldNotError)
		})
	}
}

// mergeModuleVars returns a new moduleVars, merging together all provided
// moduleVars, in order.
// Note: this is a shallow merge.
func mergeModuleVars(mvs ...moduleVars) moduleVars {
	newMvs := moduleVars{}

	for _, mv := range mvs {
		maps.Copy(newMvs, mv)
	}

	return newMvs
}
