package manager

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2Manager(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"aws_instance.runner-manager",
		"aws_internet_gateway.igw",
		"aws_route_table.rtb_public",
		"aws_route_table_association.rta_subnet_public",
		"aws_security_group.sg_22",
		"aws_subnet.subnet_public",
		"aws_vpc.vpc",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"manager override defaults": {
			moduleVars: map[string]interface{}{
				"runner_token":                           "runnerToken",
				"executor":                               "docker-autoscaler",
				"gitlab_url":                             "https://gitlab.com",
				"fleeting_service_account_access_key_id": "access_key_id",
				"fleeting_service_account_secret_access_key": "secret_access_key",
				"ssh_key_pem":           "ssh_key_pem",
				"ssh_key_pem_name":      "ssh_key_pem_name",
				"aws_asg_name":          "aws_asg_name",
				"fleeting_service":      "ec2",
				"capacity_per_instance": 1,
				"scale_min":             1,
				"scale_max":             1,
				"name":                  name,
			},
			expectedModules: expectedModules,
		},
		"default manager": {
			moduleVars: map[string]interface{}{
				"runner_token":          "runnerToken",
				"executor":              "docker-autoscaler",
				"gitlab_url":            "https://gitlab.com",
				"aws_asg_name":          "aws_asg_name",
				"fleeting_service":      "ec2",
				"capacity_per_instance": 1,
				"scale_min":             1,
				"scale_max":             1,
				"name":                  name,
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
