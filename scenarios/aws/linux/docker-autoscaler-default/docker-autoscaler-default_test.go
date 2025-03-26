package scenarios_aws_linux_docker_autoscaler_default

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestScenarioAWSLinuxDockerAutoscaler(t *testing.T) {
	name := test_tools.JobName(t)

	testCases := map[string]struct {
		moduleVars      map[string]interface{}
		expectedModules []string
	}{
		"aws linux docker autoscaler default": {
			moduleVars: map[string]interface{}{
				"name":                  name,
				"labels":                map[string]string{},
				"aws_region":            "us-east-1",
				"aws_zone":              "us-east-1b",
				"gitlab_url":            "https://gitlab.com",
				"gitlab_project_id":     "123121213",
				"capacity_per_instance": 1,
				"max_instances":         20,
				"max_use_count":         1,
				"concurrent":            1,
				"autoscaling_policy": map[string]interface{}{
					"scale_min":          1,
					"idle_time":          "2m0s",
					"scale_factor":       0,
					"scale_factor_limit": 0,
				},
				"ephemeral_runner": map[string]interface{}{
					"disk_type":    "gp3",
					"disk_size":    25,
					"machine_type": "t3.medium",
					"source_image": "",
				},
				"runner_description": "example-grit-docker-autoscaler-runner",
				"runner_tags":        []string{"grit-runner"},
			},
			expectedModules: []string{
				"module.cache.aws_iam_policy.cache_bucket_access_policy",
				"module.vpc.aws_subnet.jobs_vpc_subnet",
				"module.fleeting.module.ec2[0].module.common.aws_autoscaling_group.fleeting_asg",
				"module.cache.aws_s3_bucket.cache_bucket_server_logs",
				"module.cache.aws_s3_bucket_lifecycle_configuration.cache",
				"module.fleeting.module.ec2[0].module.common.aws_launch_template.fleeting_asg_template",
				"module.cache.aws_iam_user.cache",
				"module.cache.aws_iam_user_policy_attachment.cache_bucket_user_policy",
				"module.iam.aws_iam_user_policy_attachment.fleeting_service_account_attach",
				"module.cache.aws_s3_bucket_public_access_block.cache",
				"module.cache.aws_s3_bucket_versioning.cache",
				"module.cache.aws_s3_bucket_versioning.cache_bucket_server_logs",
				"module.gitlab.gitlab_user_runner.primary",
				"module.iam.aws_iam_policy.fleeting_service_account_policy",
				"module.vpc.aws_internet_gateway.igw",
				"module.vpc.aws_vpc.vpc",
				"module.security_groups.aws_security_group.jobs_security_group",
				"module.vpc.aws_route_table.rtb_public",
				"module.cache.aws_s3_bucket.cache",
				"module.iam.aws_iam_access_key.fleeting_service_account_key",
				"module.runner.module.ec2[0].data.cloudinit_config.config",
				"module.runner.module.ec2[0].aws_instance.runner_manager",
				"module.cache.data.aws_iam_policy_document.cache_bucket_access_policy_document",
				"module.cache.aws_iam_access_key.cache_bucket_user_key",
				"module.cache.aws_s3_bucket_logging.cache",
				"module.iam.aws_iam_user.fleeting_service_account",
				"module.fleeting.module.ec2[0].module.common.aws_key_pair.jobs_key_pair",
				"module.cache.aws_s3_bucket_public_access_block.cache_bucket_server_logs",
				"module.security_groups.aws_security_group.manager_sg",
				"module.vpc.aws_route.internet_route",
				"module.vpc.aws_route_table_association.rta_subnet_public",
				"module.fleeting.module.ec2[0].module.common.tls_private_key.aws_jobs_private_key",
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssert(t, tc.moduleVars, tc.expectedModules)
		})
	}
}
