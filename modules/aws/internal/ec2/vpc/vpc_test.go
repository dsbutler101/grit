package manager

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2Manager(t *testing.T) {
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
				"aws_vpc_cidr":        "10.0.0.0/16",
				"aws_vpc_subnet_cidr": "10.0.0.0/24",
				"aws_zone":            "us-east-1a",
				"labels":              map[string]string{"env": "another place"},
				"name":                name,
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
