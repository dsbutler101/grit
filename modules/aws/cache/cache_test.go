package cache

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestRunner(t *testing.T) {
	name := test_tools.JobName(t)

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"s3 cache bucket": {
			moduleVars: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
				"cache_object_lifetime": 10,
				"bucket_name":           "cache_bucket",
			},
			expectedModules: []string{
				"aws_s3_bucket.cache_bucket_server_logs",
				"aws_s3_bucket_versioning.cache_bucket_server_logs",
				"aws_s3_bucket_public_access_block.cache_bucket_server_logs",
				"aws_s3_bucket.cache",
				"aws_s3_bucket_versioning.cache",
				"aws_s3_bucket_public_access_block.cache",
				"aws_s3_bucket_lifecycle_configuration.cache",
				"aws_s3_bucket_logging.cache",
				"aws_iam_user.cache",
				"data.aws_iam_policy_document.cache_bucket_access_policy_document",
				"aws_iam_policy.cache_bucket_access_policy",
				"aws_iam_user_policy_attachment.cache_bucket_user_policy",
				"aws_iam_access_key.cache_bucket_user_key",
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
