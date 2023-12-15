package aws_test

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestAWSInternalEC2Fleeting(t *testing.T) {
	name := test_tools.JobName(t)

	baseExpectedModules := []string{
		"module.test-module.module.ec2-instance-group[0].module.common.aws_autoscaling_group.fleeting-asg",
		"module.test-module.module.ec2-instance-group[0].module.common.aws_iam_access_key.fleeting-service-account-key",
		"module.test-module.module.ec2-instance-group[0].module.common.aws_iam_policy.fleeting-service-account-policy",
		"module.test-module.module.ec2-instance-group[0].module.common.aws_iam_user.fleeting-service-account",
		"module.test-module.module.ec2-instance-group[0].module.common.aws_iam_user_policy_attachment.fleeting-service-account-attach",
		"module.test-module.module.ec2-instance-group[0].module.common.aws_key_pair.jobs-key-pair",
		"module.test-module.module.ec2-instance-group[0].module.common.aws_launch_template.fleeting-asg-template",
		"module.test-module.module.ec2-instance-group[0].module.common.aws_security_group.jobs-security-group",
		"module.test-module.module.ec2-instance-group[0].module.common.tls_private_key.aws-jobs-private-key",
		"module.test-module.module.ec2-managers[0].aws_instance.runner-manager",
		"module.test-module.module.ec2-managers[0].aws_security_group.manager_sg",
		"module.test-module.module.ec2-managers[0].data.cloudinit_config.config",
		"module.test-module.module.vpc[0].aws_internet_gateway.igw",
		"module.test-module.module.vpc[0].aws_route.internet-route",
		"module.test-module.module.vpc[0].aws_route_table.rtb_public",
		"module.test-module.module.vpc[0].aws_route_table_association.rta_subnet_public",
		"module.test-module.module.vpc[0].aws_subnet.jobs-vpc-subnet",
		"module.test-module.module.vpc[0].aws_vpc.vpc",
	}

	expectedModulesWithNewRunnerToken := append(baseExpectedModules, "module.test-module.module.gitlab[0].gitlab_user_runner.primary")

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"ec2-ec2-docker-autoscaler runner_token provided": {
			moduleVars: map[string]interface{}{
				"manager_service":           "ec2",
				"fleeting_service":          "ec2",
				"gitlab_project_id":         test_tools.GritEndToEndTestProjectID,
				"gitlab_runner_description": "my-linux-runner",
				"gitlab_runner_tags":        []string{"linux-ec2-job"},
				"runner_token":              "tokenString",
				"name":                      name,
			},
			expectedModules: baseExpectedModules,
		},
		"ec2-ec2-docker-autoscaler runner_token not provided": {
			moduleVars: map[string]interface{}{
				"manager_service":           "ec2",
				"fleeting_service":          "ec2",
				"gitlab_project_id":         test_tools.GritEndToEndTestProjectID,
				"gitlab_runner_description": "my-linux-runner",
				"gitlab_runner_tags":        []string{"linux-ec2-job"},
				"name":                      name,
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
