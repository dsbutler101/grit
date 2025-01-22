#######################
# METADATA VALIDATION #
#######################

module "validate-name" {
  source = "../../../internal/validation/name"
  name   = var.metadata.name
}

module "validate-support" {
  source   = "../../../internal/validation/support"
  use_case = "${var.service}-${var.executor}"
  use_case_support = tomap({
    "ec2-docker-autoscaler" = "experimental"
  })
  min_support = var.metadata.min_support
}

######################
# RUNNER PROD MODULE #
######################

module "validate-scale-parameters" {
  source                = "../../../internal/validation/scale_parameters"
  executor              = var.executor
  capacity_per_instance = var.capacity_per_instance
  scale_min             = var.scale_min
  scale_max             = var.scale_max
  idle_percentage       = var.idle_percentage
}

module "ec2" {
  count  = var.service == "ec2" ? 1 : 0
  source = "../internal/ec2"

  gitlab   = var.gitlab
  fleeting = var.fleeting
  iam      = var.iam
  vpc      = var.vpc
  s3_cache = var.s3_cache

  executor                   = var.executor
  capacity_per_instance      = var.capacity_per_instance
  scale_min                  = var.scale_min
  scale_max                  = var.scale_max
  idle_percentage            = var.idle_percentage
  max_use_count              = var.max_use_count
  security_group_ids         = var.security_group_ids
  privileged                 = var.privileged
  region                     = var.region
  runner_repository          = var.runner_repository
  runner_version             = var.runner_version
  aws_plugin_version         = var.aws_plugin_version
  instance_role_profile_name = var.instance_role_profile_name
  install_cloudwatch_agent   = var.install_cloudwatch_agent
  cloudwatch_agent_json      = var.cloudwatch_agent_json
  enable_metrics_export      = var.enable_metrics_export
  metrics_export_endpoint    = var.metrics_export_endpoint
  default_docker_image       = var.default_docker_image

  name   = var.metadata.name
  labels = var.metadata.labels
}
