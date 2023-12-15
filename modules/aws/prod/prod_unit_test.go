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
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_key_pair.jobs-key-pair",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_launch_template.fleeting-asg-template",
		"module.prod-module.module.ec2-instance-group[0].module.common.aws_security_group.jobs-security-group",
		"module.prod-module.module.ec2-instance-group[0].module.common.tls_private_key.aws-jobs-private-key",
		"module.prod-module.module.ec2-managers[0].aws_instance.runner-manager",
		"module.prod-module.module.ec2-managers[0].aws_security_group.manager_sg",
		"module.prod-module.module.ec2-managers[0].data.cloudinit_config.config",
	}

	baseModulesPlusNewRunnerToken := append(baseExpectedModules, "module.prod-module.module.gitlab[0].gitlab_user_runner.primary")

	vpcExpectedModules := []string{
		"module.prod-module.module.vpc[0].aws_internet_gateway.igw",
		"module.prod-module.module.vpc[0].aws_route.internet-route",
		"module.prod-module.module.vpc[0].aws_route_table.rtb_public",
		"module.prod-module.module.vpc[0].aws_route_table_association.rta_subnet_public",
		"module.prod-module.module.vpc[0].aws_subnet.jobs-vpc-subnet",
		"module.prod-module.module.vpc[0].aws_vpc.vpc",
	}

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"ec2-ec2-docker-autoscaler runner_token provided and preexisting vpc": {
			moduleVars: map[string]interface{}{
				"gitlab_project_id":         test_tools.GritEndToEndTestProjectID,
				"gitlab_runner_description": "my-linux-runner",
				"gitlab_runner_tags":        []string{"linux-ec2-job"},
				"runner_token":              "tokenString",
				"manager_service":           "ec2",
				"fleeting_service":          "ec2",
				"fleeting_os":               "linux",
				"executor":                  "docker-autoscaler",
				"scale_min":                 2,
				"scale_max":                 10,
				"idle_percentage":           10,
				"capacity_per_instance":     1,
				"ami":                       "ami-05012401516a40259",
				"instance_type":             "t2.medium",
				"aws_vpc_id":                "12345",
				"aws_vpc_subnet_id":         "12345",
				"min_maturity":              "alpha",
				"name":                      name,
			},
			expectedModules: baseExpectedModules,
		},
		"ec2-ec2-docker-autoscaler runner_token provided": {
			moduleVars: map[string]interface{}{
				"gitlab_project_id":         test_tools.GritEndToEndTestProjectID,
				"gitlab_runner_description": "my-linux-runner",
				"gitlab_runner_tags":        []string{"linux-ec2-job"},
				"runner_token":              "tokenString",
				"manager_service":           "ec2",
				"fleeting_service":          "ec2",
				"fleeting_os":               "linux",
				"executor":                  "docker-autoscaler",
				"scale_min":                 2,
				"scale_max":                 10,
				"idle_percentage":           10,
				"capacity_per_instance":     1,
				"ami":                       "ami-05012401516a40259",
				"instance_type":             "t2.medium",
				"aws_vpc_cidr":              "10.0.0.0/16",
				"aws_vpc_subnet_cidr":       "10.0.0.0/24",
				"min_maturity":              "alpha",
				"name":                      name,
			},
			expectedModules: append(baseExpectedModules, vpcExpectedModules...),
		},
		"ec2-ec2-docker-autoscaler runner_token not provided": {
			moduleVars: map[string]interface{}{
				"gitlab_project_id":         test_tools.GritEndToEndTestProjectID,
				"gitlab_runner_description": "my-linux-runner",
				"gitlab_runner_tags":        []string{"linux-ec2-job"},
				"manager_service":           "ec2",
				"fleeting_service":          "ec2",
				"fleeting_os":               "linux",
				"executor":                  "docker-autoscaler",
				"scale_min":                 2,
				"scale_max":                 10,
				"idle_percentage":           10,
				"capacity_per_instance":     1,
				"ami":                       "ami-05012401516a40259",
				"instance_type":             "t2.medium",
				"aws_vpc_cidr":              "10.0.0.0/16",
				"aws_vpc_subnet_cidr":       "10.0.0.0/24",
				"min_maturity":              "alpha",
				"name":                      name,
			},
			expectedModules: append(baseModulesPlusNewRunnerToken, vpcExpectedModules...),
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
