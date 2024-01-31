package manager

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestRunner(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"google_project_iam_custom_role.runner-manager",
		"google_project_iam_member.runner-manager",
		"google_compute_firewall.runner-manager-ssh-access",
		"data.cloudinit_config.config",
		"google_compute_instance.runner-manager",
		"google_kms_key_ring.default",
		"google_kms_crypto_key.default",
		"google_kms_secret_ciphertext.runner-token",
		"google_kms_secret_ciphertext.runner-ssh-key",
		"tls_private_key.runner-manager",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"create runner": {
			moduleVars: map[string]interface{}{
				"name":                   name,
				"labels":                 map[string]string{"env": "another place"},
				"google_project":         "example-project-a1b2c3",
				"runner_version":         "v16.8.0",
				"machine_type":           "e2-micro",
				"disk_size_gb":           25,
				"disk_type":              "pd-standard",
				"service_account_email":  "service-account@example.com",
				"concurrent":             5,
				"check_interval":         10,
				"log_level":              "info",
				"listen_address":         ":9402",
				"gitlab_url":             "https://gitlab.example.com",
				"runner_token":           "glrt-SOME_TOKEN",
				"request_concurrency":    5,
				"cache_gcs_bucket":       "cache-bucket",
				"runners_global_section": "global section custom part",
				"runners_docker_section": "docker section custom part",
				"vpc": map[string]string{
					"id":        "vpc-id",
					"subnet_id": "subnet-id",
				}},
			expectedModules: expectedModules,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
