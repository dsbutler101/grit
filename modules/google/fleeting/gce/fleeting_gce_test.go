package manager

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestFleetingGCE(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"google_project_iam_custom_role.instance_group_manager",
		"google_project_iam_member.instance_group_manager",
		"google_compute_instance_template.ephemeral_runner",
		"google_compute_instance_group_manager.ephemeral_runners",
		"google_compute_firewall.ephemeral_runners_ssh_access",
		"google_compute_firewall.ephemeral_runners_cross_vm_deny",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"create fleeting gce": {
			moduleVars: map[string]interface{}{
				"name":                  name,
				"labels":                map[string]string{"env": "another-place"},
				"google_project":        "example-project-a1b2c3",
				"subnetwork_project":    "example-subnetwork-project-a1b2c3",
				"google_zone":           "us-central1-a",
				"service_account_email": "service-account@example.com",
				"machine_type":          "e2-micro",
				"disk_type":             "pd-standard",
				"disk_size_gb":          25,
				"source_image":          "some-source-image",
				"vpc": map[string]interface{}{
					"enabled": true,
					"id":      "my-vpc",
					"subnetwork_ids": map[string]any{
						"runner-manager":    "manager-subnet",
						"ephemeral-runners": "runners-subnet",
					},
					"subnetwork_cidrs": map[string]any{
						"runner-manager":    "manager-subnet",
						"ephemeral-runners": "runners-subnet",
					},
				},
				"manager_subnet_name": "runner-manager",
				"runners_subnet_name": "ephemeral-runners",
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
