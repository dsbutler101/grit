package prometheus

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestPrometheus(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"google_compute_instance.prometheus_server",
		"google_compute_firewall.prometheus_ssh_access",
		"terraform_data.prometheus_server_replacement",
		"google_compute_disk.prometheus_data",
		"google_project_iam_custom_role.prometheus_server",
		"google_project_iam_member.prometheus_server",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"create prometheus": {
			moduleVars: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{"env": "another-place"},
					"min_support": "experimental",
				},
				"google_project":        "example-project-a1b2c3",
				"google_zone":           "us-east1-c",
				"service_account_email": "service-account@example.com",
				"machine_type":          "example-machine-type",
				"boot_disk": map[string]any{
					"disk_type": "pd-ssd",
					"size_gb":   25,
				},
				"data_disk": map[string]any{
					"disk_type": "pd-ssd",
					"size_gb":   100,
				},
				"prometheus_version":    "latet",
				"node_exporter_version": "latest",
				"node_exporter_port":    9100,
				"prometheus_external_labels": map[string]string{
					"label_1": "value_1",
					"label_2": "value_2",
				},
				"mimir": map[string]any{
					"url":    "https://mimir.example.com",
					"tenant": "example-tenant",
				},
				"runner_manager_nodes": map[string]any{
					"filter": "some-filter",
					"exporter_ports": map[string]int{
						"runner_manager": 9402,
						"node_exporter":  9100,
					},
					"custom_relabel_configs": []map[string]any{
						{
							"target_label":  "some_label",
							"source_labels": []string{"other_label_1", "other_label_2"},
							"regex":         "(.*)",
							"replacement":   "$1",
							"action":        "replace",
						},
					},
					"instance_labels_to_include": []string{"label_1", "label_2"},
				},
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
