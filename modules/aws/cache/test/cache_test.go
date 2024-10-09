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
		"s3 cahe bucket": {
			moduleVars: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":   name,
					"labels": map[string]string{},
				},
				"cache_object_lifetime": 10,
				"bucket_name":           "cache-bucket",
			},
			expectedModules: []string{
				"module.cache.aws_s3_bucket.cache-bucket-server-logs",
				"module.cache.aws_s3_bucket_versioning.cache-bucket-server-logs",
				"module.cache.aws_s3_bucket_public_access_block.cache-bucket-server-logs",
				"module.cache.aws_s3_bucket.cache",
				"module.cache.aws_s3_bucket_versioning.cache",
				"module.cache.aws_s3_bucket_public_access_block.cache",
				"module.cache.aws_s3_bucket_lifecycle_configuration.cache",
				"module.cache.aws_s3_bucket_logging.cache",
				"module.cache.aws_iam_user.cache",
				"module.cache.data.aws_iam_policy_document.cache-bucket-access-policy-document",
				"module.cache.aws_iam_policy.cache-bucket-access-policy",
				"module.cache.aws_iam_user_policy_attachment.cache-bucket-user-policy",
				"module.cache.aws_iam_access_key.cache-bucket-user-key",
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
