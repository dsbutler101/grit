package aws_prod

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSProd(t *testing.T) {
	name := test_tools.JobName(t)

	baseExpectedModules := []string{
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_autoscaling_group.fleeting-asg",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_iam_access_key.fleeting-service-account-key",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_iam_policy.fleeting-service-account-policy",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_iam_user.fleeting-service-account",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_iam_user_policy_attachment.fleeting-service-account-attach",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_internet_gateway.internet-access",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_key_pair.jobs-key-pair",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_launch_template.fleeting-asg-template",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_route.internet-route",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_security_group.jobs-security-group",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_subnet.jobs-vpc-subnet",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_vpc.jobs-vpc",
		"module.prod-module.module.ec2-instance-group[0].module.common.data.aws_route_table.jobs-route-table",
		"module.prod-module.module.ec2-instance-group[0].module.common.tls_private_key.aws-jobs-private-key",
		"module.prod-module.module.ec2-managers[0].aws_instance.runner-manager",
		"module.prod-module.module.ec2-managers[0].aws_internet_gateway.igw",
		"module.prod-module.module.ec2-managers[0].aws_route_table.rtb_public",
		"module.prod-module.module.ec2-managers[0].aws_route_table_association.rta_subnet_public",
		"module.prod-module.module.ec2-managers[0].aws_security_group.sg_22",
		"module.prod-module.module.ec2-managers[0].aws_subnet.subnet_public",
		"module.prod-module.module.ec2-managers[0].aws_vpc.vpc",
		"module.prod-module.module.ec2-managers[0].data.cloudinit_config.config",
	}

	expectedModulesWithNewRunnerToken := append(baseExpectedModules, "module.prod-module.module.gitlab[0].gitlab_user_runner.primary")

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"ec2-ec2-docker-autoscaler runner_token provided": {
			moduleVars: map[string]interface{}{
				"manager_service":           "ec2",
				"fleeting_service":          "ec2",
				"fleeting_os":               "linux",
				"ami":                       "ami-05012401516a40259",
				"executor":                  "docker-autoscaler",
				"instance_type":             "t2.medium",
				"aws_vpc_cidr":              "10.0.0.0/24",
				"capacity_per_instance":     1,
				"scale_min":                 2,
				"scale_max":                 10,
				"gitlab_project_id":         test_tools.GritEndToEndTestProjectID,
				"gitlab_runner_description": "my-linux-runner",
				"gitlab_runner_tags":        []string{"linux-ec2-job"},
				"runner_token":              "tokenString",
				"name":                      name,
				"min_maturity":              "alpha",
			},
			expectedModules: baseExpectedModules,
		},
		"ec2-ec2-docker-autoscaler runner_token not provided": {
			moduleVars: map[string]interface{}{
				"manager_service":           "ec2",
				"fleeting_service":          "ec2",
				"fleeting_os":               "linux",
				"ami":                       "ami-05012401516a40259",
				"executor":                  "docker-autoscaler",
				"instance_type":             "t2.medium",
				"aws_vpc_cidr":              "10.0.0.0/24",
				"capacity_per_instance":     1,
				"scale_min":                 2,
				"scale_max":                 10,
				"gitlab_project_id":         test_tools.GritEndToEndTestProjectID,
				"gitlab_runner_description": "my-linux-runner",
				"gitlab_runner_tags":        []string{"linux-ec2-job"},
				"name":                      name,
				"min_maturity":              "alpha",
			},
			expectedModules: expectedModulesWithNewRunnerToken,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
