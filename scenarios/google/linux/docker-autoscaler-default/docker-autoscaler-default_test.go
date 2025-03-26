package scenarios_google_linux_docker_autoscaler_default

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestScenarioGoogleLinuxDockerAutoscaler(t *testing.T) {
	name := test_tools.JobName(t)

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"google linux docker autoscaler default": {
			moduleVars: map[string]interface{}{
				"name":           name,
				"labels":         map[string]string{"env": "prod", "managed-by": "terraform"},
				"google_project": "my-gcp-project",
				"google_region":  "us-central1",
				"google_zone":    "us-central1-a",
				"gitlab_url":     "https://gitlab.example.com",
				"runner_token":   "glrt-abcdef1234567890",

				"runner_machine_type": "n2d-standard-4",
				"runner_disk_type":    "pd-ssd",
				"concurrent":          50,

				"runners_global_section": "",
				"runners_docker_section": "",
				"capacity_per_instance":  1,
				"max_instances":          200,
				"max_use_count":          1,

				"autoscaling_policies": []map[string]interface{}{
					{
						"periods":            []string{"* * * * *"},
						"timezone":           "UTC",
						"scale_min":          3,
						"idle_time":          "20m0s",
						"scale_factor":       0,
						"scale_factor_limit": 0,
					},
					{
						"periods":            []string{"0 9-18 * * 1-5"},
						"timezone":           "",
						"scale_min":          10,
						"idle_time":          "10m0s",
						"scale_factor":       1.5,
						"scale_factor_limit": 20,
					},
				},

				"ephemeral_runner": map[string]interface{}{
					"disk_type":    "pd-standard",
					"disk_size":    25,
					"machine_type": "n2d-standard-2",
					"source_image": "projects/cos-cloud/global/images/family/cos-stable",
				},

				"prometheus": map[string]interface{}{
					"enabled": true,
					"mimir": map[string]interface{}{
						"url":    "https://mimir.example.com/api/v1/push",
						"tenant": "gitlab-runner",
					},
					"external_labels": map[string]string{
						"environment": "production",
						"service":     "gitlab-runner",
					},
					"custom_relabel_configs": []map[string]interface{}{
						{
							"target_label":  "job",
							"source_labels": []string{""},
							"regex":         "(.*)",
							"replacement":   "$1",
							"action":        "replace",
						},
					},
					"instance_labels_to_include": []string{"instance", "job"},
				},
			},
			expectedModules: []string{
				"module.runner.google_kms_secret_ciphertext.runner_ssh_key",
				"module.runner.google_kms_crypto_key.default",
				"module.runner.google_project_iam_custom_role.runner_manager",
				"module.fleeting.module.gce[0].google_compute_instance_group_manager.ephemeral_runners",
				"module.fleeting.module.gce[0].google_project_iam_member.instance_group_manager",
				"module.cache.google_storage_bucket.cache_bucket",
				"module.prometheus[0].google_project_iam_member.prometheus_server",
				"module.prometheus_iam.google_service_account.default",
				"module.runner.google_compute_firewall.runner_manager_ssh_access",
				"module.fleeting.module.gce[0].google_compute_instance_template.ephemeral_runner",
				"module.cache.google_storage_bucket_iam_binding.cache_bucket",
				"module.runner.google_compute_instance.runner_manager",
				"module.runner.google_project_iam_member.runner_manager",
				"module.prometheus[0].google_compute_disk.prometheus_data",
				"module.vpc.google_compute_firewall.runner_manager_ingress_default",
				"module.runner.tls_private_key.runner_manager",
				"module.vpc.google_compute_firewall.runner_manager_egress_default",
				"module.prometheus[0].google_compute_firewall.prometheus_ssh_access",
				"module.prometheus[0].google_project_iam_custom_role.prometheus_server",
				"module.prometheus[0].terraform_data.prometheus_server_replacement",
				"module.runner.data.cloudinit_config.config",
				"module.runner.google_compute_address.runner_manager",
				"module.runner.google_kms_key_ring.default",
				"module.fleeting.module.gce[0].google_compute_firewall.ephemeral_runners_ssh_access",
				"module.iam.google_service_account.default",
				"module.runner.google_kms_secret_ciphertext.runner_token",
				"module.vpc.google_compute_subnetwork.subnetwork[\"ephemeral-runners\"]",
				"module.vpc.google_compute_subnetwork.subnetwork[\"runner-manager\"]",
				"module.fleeting.module.gce[0].google_compute_firewall.ephemeral_runners_cross_vm_deny",
				"module.prometheus[0].google_compute_instance.prometheus_server",
				"module.runner.google_compute_firewall.additional_rules[\"prometheus\"]",
				"module.vpc.google_compute_network.default",
				"module.vpc.google_compute_subnetwork.subnetwork[\"prometheus\"]",
				"module.fleeting.module.gce[0].google_project_iam_custom_role.instance_group_manager",
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
