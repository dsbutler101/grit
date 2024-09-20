package manager

import (
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

type moduleVars = map[string]any

var defaultModuleVars = moduleVars{
	"name":          test_tools.JobName(nil),
	"google_region": "us-east1",
	"google_zone":   "us-east1-b",
	"labels": map[string]string{
		"some": "label",
	},
	"node_pools": map[string]any{
		"default": map[string]any{
			"node_count": 1,
			"node_config": map[string]any{
				"machine_type": "e2-micro",
			},
		},
	},
	"vpc": map[string]string{
		"id":        "",
		"subnet_id": "",
	},
	"deletion_protection": "false",
}

func TestGKE(t *testing.T) {
	expectedModules := []string{
		"google_container_cluster.primary",
		`google_container_node_pool.node_pool["default"]`,
	}

	testCases := map[string]struct {
		moduleVars      moduleVars
		expectedModules []string
	}{
		"create-gke-happy": {
			moduleVars: moduleVars{
				"labels": map[string]string{"env": "another-place"},
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
			mvs := mergeModuleVars(
				defaultModuleVars,
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
				"labels": map[string]string{
					"label1": "google does not like this because it has spaces",
					"label2": "this-is-fine",
				},
			},
		},
		"valid-label-value": {
			moduleVars: moduleVars{
				"labels": map[string]string{
					"label1": "this-is-fine",
					"label2": "this.is_fine-too",
				},
			},
			shouldNotError: true,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			mvs := mergeModuleVars(
				defaultModuleVars,
				tc.moduleVars,
			)
			test_tools.PlanAndAssertError(t, mvs, !tc.shouldNotError)
		})
	}
}

// TestGKEVarsUnused is to test planning with module vars that are either not
// used yet or anymore (deprecated)
// When unused module variables are changed across two planning runs, we expect
// there to be
//   - no diff between the two plans
//   - no error on either of the planning runs
func TestGKEVarsUnused(t *testing.T) {
	testCases := map[string]struct {
		sharedVars    moduleVars
		varsRun1      moduleVars
		varsRun2      moduleVars
		expectChanges bool
	}{
		"google-region-is-not-used": {
			varsRun1: moduleVars{"google_region": "some-region"},
			varsRun2: moduleVars{"google_region": "some-other-region"},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			mv1 := mergeModuleVars(defaultModuleVars, tc.sharedVars, tc.varsRun1)
			mv2 := mergeModuleVars(defaultModuleVars, tc.sharedVars, tc.varsRun2)

			resourceDiff1 := test_tools.Plan(t, mv1).RawPlan.ResourceChanges
			resourceDiff2 := test_tools.Plan(t, mv2).RawPlan.ResourceChanges

			if !tc.expectChanges {
				assert.Equal(t, resourceDiff1, resourceDiff2, "Expected plans not to differ")
			} else {
				assert.NotEqual(t, resourceDiff1, resourceDiff2, "Expected plans to differ, but they don't")
			}
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
