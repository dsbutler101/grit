package securitygroups

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

type SGRule struct {
	FromPort    int      `json:"from_port"`
	ToPort      int      `json:"to_port"`
	Protocol    string   `json:"protocol"`
	CidrBlocks  []string `json:"cidr_blocks"`
	Description string   `json:"description"`
}

func TestAWSSecurityGroups(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"aws_security_group.manager_sg",
		"aws_security_group.jobs_security_group",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"common fleet": {
			moduleVars: map[string]interface{}{
				"labels": map[string]string{"env": "another place"},
				"vpc_id": "1234",
				"name":   name,
				"manager_outbound_sg_rules": map[string]interface{}{
					"outbound_all": map[string]interface{}{
						"from_port":   0,
						"to_port":     0,
						"protocol":    "-1",
						"cidr_blocks": []string{"0.0.0.0/0"},
						"description": "allow outbound",
					},
				},
				"fleeting_outbound_sg_rules": map[string]interface{}{
					"outbound_all": map[string]interface{}{
						"from_port":   0,
						"to_port":     0,
						"protocol":    "-1",
						"cidr_blocks": []string{"0.0.0.0/0"},
						"description": "allow outbound",
					},
				},
				"manager_inbound_sg_rules": map[string]interface{}{
					"inbound_all": map[string]interface{}{
						"from_port":   0,
						"to_port":     0,
						"protocol":    "-1",
						"cidr_blocks": []string{"0.0.0.0/0"},
						"description": "allow inbound",
					},
				},
				"fleeting_inbound_sg_rules": map[string]interface{}{
					"inbound_all": map[string]interface{}{
						"from_port":   0,
						"to_port":     0,
						"protocol":    "-1",
						"cidr_blocks": []string{"0.0.0.0/0"},
						"description": "allow inbound",
					},
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
