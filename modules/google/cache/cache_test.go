package cache

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestCache(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"google_storage_bucket.cache-bucket",
		"google_storage_bucket_iam_binding.cache-bucket",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"create cache": {
			moduleVars: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{"env": "another place"},
					"min_support": "experimental",
				},
				"cache_object_lifetime": 10,
				"bucket_name":           "optional-explicit-bucket-name",
				"bucket_location":       "us-east1",
				"service_account_emails": []string{
					"service-account@example.com",
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
