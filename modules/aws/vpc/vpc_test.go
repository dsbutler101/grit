package vpc

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestVPC(t *testing.T) {
	name := test_tools.JobName(t)

	expectedModules := []string{
		"aws_internet_gateway.igw",
		"aws_route.internet-route",
		"aws_route_table.rtb_public",
		"aws_route_table_association.rta_subnet_public",
		"aws_subnet.jobs-vpc-subnet",
		"aws_vpc.vpc",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"create vpc": {
			moduleVars: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":        name,
					"labels":      map[string]string{"env": "another place"},
					"min_support": "experimental",
				},
				"cidr":        "10.0.0.0/16",
				"subnet_cidr": "10.0.0.0/24",
				"zone":        "us-east-1a",
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
