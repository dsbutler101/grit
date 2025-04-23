package runner

import (
	"maps"
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestRunner(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"google_project_iam_custom_role.runner_manager",
		"google_project_iam_member.runner_manager",
		"google_compute_firewall.runner_manager_ssh_access",
		"data.cloudinit_config.config",
		"google_compute_instance.runner_manager",
		"google_compute_address.runner_manager",
		"google_kms_key_ring.default",
		"google_kms_crypto_key.default",
		"google_kms_secret_ciphertext.runner_token",
		"google_kms_secret_ciphertext.runner_ssh_key",
		"tls_private_key.runner_manager",
	}

	variables := func(overrides map[string]any) map[string]any {
		required := map[string]any{
			"metadata": map[string]interface{}{
				"name":        name,
				"labels":      map[string]string{"env": "another-place"},
				"min_support": "experimental",
			},
			"google_project":               "example-project-a1b2c3",
			"subnetwork_project":           "example-subnetwork-project-a1b2c3",
			"address_type":                 "EXTERNAL",
			"access_config_enabled":        true,
			"google_zone":                  "us-east1-b",
			"service_account_email":        "service-account@example.com",
			"gitlab_url":                   "https://gitlab.example.com",
			"runner_token":                 "glrt-SOME_TOKEN",
			"executor":                     "docker-autoscaler",
			"fleeting_instance_group_name": "instance-group-name",
			"vpc": map[string]any{
				"enabled": true,
				"id":      "vpc-id",
				"subnetwork_ids": map[string]any{
					"runner-manager": "subnet-id",
				},
				"subnetwork_cidrs": map[string]any{
					"runner-manager": "subnet-cidr",
				},
			},
		}
		maps.Copy(required, overrides)
		return required
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"with required variables": {
			moduleVars:      variables(nil),
			expectedModules: expectedModules,
		},
		"with optional variables": {
			moduleVars: variables(map[string]any{
				"node_exporter": map[string]any{
					"version": "latest",
					"port":    9100,
				},
				"machine_type":   "e2-micro",
				"disk_image":     "image",
				"disk_type":      "disk",
				"disk_size_gb":   25,
				"runner_version": "v17.8.0",
				"concurrent":     6,
				"check_interval": 5,
				"log_level":      "warn",
				"listen_address": ":9402",
				"runner_metrics_listener": map[string]any{
					"address": "1.2.3.4",
					"port":    9402,
				},
				"request_concurrency":                   2,
				"cache_gcs_bucket":                      "cache-bucket",
				"runners_global_section":                "global section\ncustom part",
				"runners_docker_section":                "docker section\ncustom part",
				"default_docker_image":                  "alpine:latest",
				"fleeting_googlecompute_plugin_version": "v0.0.1",
				"capacity_per_instance":                 5,
				"max_instances":                         5,
				"max_use_count":                         5,
				"autoscaling_policies": []map[string]any{
					{
						"periods":            []string{"* * * * *"},
						"timezone":           "UTC",
						"scale_min":          1,
						"idle_time":          "20m0s",
						"scale_factor":       0.5,
						"scale_factor_limit": 10,
					},
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
				"source_ranges":         []string{"0.0.0.0/8"},
				"kms_location":          "local",
				"address_type":          "INTERNAL",
				"access_config_enabled": false,
				"additional_tags":       []string{"potato"},
				"runner_registry":       "gitlab-runner",
				"https_proxy":           "localhost:1234",
				"no_proxy":              "foo",
				"additional_volumes":    []string{"a", "b"},
				"runner_wrapper": map[string]any{
					"enabled":                     true,
					"process_termination_timeout": "1h",
					"socket_path":                 "tcp://localhost:1234",
				},
				"vpc": map[string]any{
					"enabled": true,
					"id":      "vpc-id",
					"subnetwork_ids": map[string]any{
						"runner-manager": "subnet-id",
					},
					"subnetwork_cidrs": map[string]any{
						"runner-manager": "subnet-cidr",
					},
				},
			}),
			expectedModules: append(expectedModules,
				`google_compute_firewall.additional_rules["test-allow"]`,
				`google_compute_firewall.additional_rules["test-deny"]`),
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
