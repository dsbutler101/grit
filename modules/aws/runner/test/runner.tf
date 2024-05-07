#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

######################
# RUNNER TEST MODULE #
######################

module "ec2" {
  count  = var.service == "ec2" ? 1 : 0
  source = "../internal/ec2"

  gitlab   = var.gitlab
  fleeting = var.fleeting
  iam      = var.iam
  vpc      = var.vpc

  executor                   = var.executor
  capacity_per_instance      = var.capacity_per_instance
  scale_min                  = var.scale_min
  scale_max                  = var.scale_max
  security_group_ids         = var.security_group_ids
  privileged                 = var.privileged
  region                     = var.region
  runner_version             = var.runner_version
  aws_plugin_version         = var.aws_plugin_version
  instance_role_profile_name = var.instance_role_profile_name
  install_cloudwatch_agent   = var.install_cloudwatch_agent
  cloudwatch_agent_json      = var.cloudwatch_agent_json

  name   = var.metadata.name
  labels = var.metadata.labels
}
