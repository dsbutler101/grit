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
					"name":        name,
					"labels":      map[string]string{},
					"min_support": "experimental",
				},
				"cache_object_lifetime": 10,
				"bucket_name":           "cache-bucket",
			},
			expectedModules: []string{
				"aws_s3_bucket.cache-bucket-server-logs",
				"aws_s3_bucket_versioning.cache-bucket-server-logs",
				"aws_s3_bucket_public_access_block.cache-bucket-server-logs",
				"aws_s3_bucket.cache",
				"aws_s3_bucket_versioning.cache",
				"aws_s3_bucket_public_access_block.cache",
				"aws_s3_bucket_lifecycle_configuration.cache",
				"aws_s3_bucket_logging.cache",
				"aws_iam_user.cache",
				"data.aws_iam_policy_document.cache-bucket-access-policy-document",
				"aws_iam_policy.cache-bucket-access-policy",
				"aws_iam_user_policy_attachment.cache-bucket-user-policy",
				"aws_iam_access_key.cache-bucket-user-key",
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
