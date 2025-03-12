package runner

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
		`google_compute_firewall.additional-rules["test-allow"]`,
		`google_compute_firewall.additional-rules["test-deny"]`,
		"data.cloudinit_config.config",
		"google_compute_instance.runner-manager",
		"google_compute_address.runner-manager",
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
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{"env": "another place"},
					"min_support": "experimental",
				},
				"google_project":        "example-project-a1b2c3",
				"subnetwork_project":    "example-subnetwork-project-a1b2c3",
				"address_type":          "EXTERNAL",
				"access_config_enabled": true,
				"google_zone":           "us-east1-b",
				"runner_version":        "v16.8.0",
				"machine_type":          "e2-micro",
				"disk_image":            "image",
				"disk_size_gb":          25,
				"disk_type":             "pd-standard",
				"service_account_email": "service-account@example.com",
				"concurrent":            5,
				"check_interval":        10,
				"log_level":             "info",
				"listen_address":        ":9402",
				"runner_metrics_listener": map[string]any{
					"address": "0.0.0.0",
					"port":    9402,
				},
				"gitlab_url":                            "https://gitlab.example.com",
				"runner_token":                          "glrt-SOME_TOKEN",
				"request_concurrency":                   5,
				"executor":                              "docker-autoscaler",
				"cache_gcs_bucket":                      "cache-bucket",
				"runners_global_section":                "global section custom part",
				"runners_docker_section":                "docker section custom part",
				"default_docker_image":                  "ubuntu:latest",
				"fleeting_googlecompute_plugin_version": "v0.1.0",
				"fleeting_instance_group_name":          "instance-group-name",
				"capacity_per_instance":                 5,
				"max_instances":                         5,
				"max_use_count":                         5,
				"autoscaling_policies": []map[string]interface{}{
					{
						"periods":            []string{"* * * * *"},
						"timezone":           "UTC",
						"scale_min":          1,
						"idle_time":          "20m0s",
						"scale_factor":       0.5,
						"scale_factor_limit": 10,
					},
				},
				"node_exporter": map[string]any{
					"version": "latest",
					"port":    9100,
				},
				"runner_manager_additional_firewall_rules": map[string]any{
					"test-allow": map[string]any{
						"direction": "INGRESS",
						"priority":  1000,
						"allow": []map[string]any{
							{
								"protocol": "tcp",
								"ports":    []int{22},
							},
						},
						"source_ranges": []string{
							"0.0.0.0/0",
						},
					},
					"test-deny": map[string]any{
						"direction": "INGRESS",
						"priority":  1000,
						"deny": []map[string]any{
							{
								"protocol": "tcp",
								"ports":    []int{53},
							},
							{
								"protocol": "udp",
								"ports":    []int{53},
							},
						},
						"source_ranges": []string{
							"0.0.0.0/0",
						},
					},
				},
				"vpc": map[string]string{
					"id":        "vpc-id",
					"subnet_id": "subnet-id",
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
