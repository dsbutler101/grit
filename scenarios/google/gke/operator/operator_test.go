package scenarios_google_gke_operator

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

var defaultModuleVars = map[string]any{
	"name":                "runner-cluster",
	"deletion_protection": true,
	"google_region":       "us-central1",
	"google_zone":         "us-central1-a",
	"subnet_cidr":         "10.0.0.0/10",
	"labels": map[string]string{
		"environment": "production",
		"managed-by":  "terraform",
		"team":        "devops",
	},
	"autoscaling": map[string]any{
		"enabled":                     true,
		"auto_provisioning_locations": []string{"us-central1-a", "us-central1-b", "us-central1-c"},
		"autoscaling_profile":         "BALANCED",
		"resource_limits": []map[string]any{
			{
				"resource_type": "cpu",
				"minimum":       1,
				"maximum":       16,
			},
			{
				"resource_type": "memory",
				"minimum":       2,
				"maximum":       32,
			},
		},
	},
	"node_pools": map[string]any{
		"default": map[string]any{
			"node_count": 2,
			"autoscaling": map[string]any{
				"min_node_count": 1,
				"max_node_count": 5,
			},
		},
		"cpu-intensive": map[string]any{
			"autoscaling": map[string]any{
				"min_node_count": 0,
				"max_node_count": 3,
			},
			"node_config": map[string]any{
				"labels": map[string]string{
					"role": "cpu-intensive",
				},
				"tags":         []string{"gitlab-runner", "cpu-intensive"},
				"machine_type": "c2-standard-8",
				"disk_size_gb": 150,
				"disk_type":    "pd-ssd",
			},
		},
	},
	"runners": map[string]any{
		"main": map[string]any{
			"runner_token":    "glrt-test-token-123456",
			"url":             "https://gitlab.com",
			"concurrent":      10,
			"check_interval":  3,
			"locked":          true,
			"protected":       true,
			"runner_tags":     []string{"gke", "kubernetes", "cloud"},
			"run_untagged":    false,
			"config_template": "",
			"envvars": map[string]string{
				"RUNNER_EXECUTOR":                 "kubernetes",
				"KUBERNETES_PRIVILEGED":           "true",
				"KUBERNETES_CPU_LIMIT":            "1",
				"KUBERNETES_MEMORY_LIMIT":         "1Gi",
				"KUBERNETES_CPU_REQUEST":          "500m",
				"KUBERNETES_MEMORY_REQUEST":       "500Mi",
				"KUBERNETES_SERVICE_ACCOUNT":      "gitlab-runner",
				"KUBERNETES_SERVICE_CPU_LIMIT":    "1",
				"KUBERNETES_SERVICE_MEMORY_LIMIT": "1Gi",
			},
			"runner_image":   "gitlab/gitlab-runner:alpine-v16.2.0",
			"helper_image":   "gitlab/gitlab-runner-helper:x86_64-v16.2.0",
			"pod_spec_patches": []map[string]any{},
			"log_level":        "info",
			"listen_address":   "[::]:9252",
		},
	},
	"operator": map[string]any{
		"version":            "v1.5.0",
		"override_manifests": "file://testdata/override.yaml",
	},
}

func TestGKE(t *testing.T) {
	expectedModules := []string{
		"module.cluster.google_container_node_pool.linux_node_pool[\"default\"]",
		"module.operator.data.http.manifest",
		"module.runner[\"main\"].terraform_data.config_template",
		"module.runner[\"main\"].terraform_data.token_secret",
		"module.operator.module.latest_operator_version.data.http.gitlab_tags",
		"module.cluster.data.google_container_engine_versions.gke_version",
		"module.cluster.google_container_node_pool.linux_node_pool[\"cpu-intensive\"]",
		"module.runner[\"main\"].kubectl_manifest.token_secret",
		"module.vpc.google_compute_network.default",
		"module.vpc.google_compute_subnetwork.subnetwork[\"runner-cluster\"]",
		"module.operator.kubectl_manifest.operator_resources[\"Namespace:some-namespace\"]",
		"module.runner[\"main\"].kubectl_manifest.envvars[0]",
		"module.runner[\"main\"].kubectl_manifest.manifest",
		"module.runner[\"main\"].terraform_data.envvars",
		"module.cluster.data.google_client_config.provider",
		"module.cluster.google_container_cluster.primary",
		"module.vpc.google_compute_firewall.runner_manager_egress_default",
		"module.vpc.google_compute_firewall.runner_manager_ingress_default",
	}

	testCases := map[string]struct {
		expectedModules []string
	}{
		"create gke operator scenario": {
			expectedModules: expectedModules,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, defaultModuleVars, tc.expectedModules)
		})
	}
}
