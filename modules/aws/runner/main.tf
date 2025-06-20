#######################
# METADATA VALIDATION #
#######################

module "validate_name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate_support" {
  source   = "../../internal/validation/support"
  use_case = "${var.service}-${var.executor}"
  use_case_support = tomap({
    "ec2-docker-autoscaler" = "experimental"
  })
  min_support = var.metadata.min_support
}

##################
# DEFAULT LABELS #
##################

module "labels" {
  source = "../../internal/labels"

  name              = var.metadata.name
  additional_labels = var.metadata.labels
}

##################
# RUNNER VERSION #
##################

module "default_version" {
  source = "../../gitlab/runner_version_lookup"

  metadata = var.metadata
  skew     = 0
}

######################
# RUNNER PROD MODULE #
######################

module "validate_scale_parameters" {
  source                = "../../internal/validation/scale_parameters"
  executor              = var.executor
  capacity_per_instance = var.capacity_per_instance
  scale_min             = var.scale_min
  scale_max             = var.scale_max
  idle_percentage       = var.idle_percentage
}

module "ec2" {
  count  = var.service == "ec2" ? 1 : 0
  source = "./ec2"

  gitlab   = var.gitlab
  fleeting = var.fleeting
  iam      = var.iam
  vpc      = var.vpc
  s3_cache = var.cache

  executor                   = var.executor
  runners_global_section     = var.runners_global_section
  capacity_per_instance      = var.capacity_per_instance
  scale_min                  = var.scale_min
  scale_max                  = var.scale_max
  idle_percentage            = var.idle_percentage
  idle_time                  = var.idle_time
  max_use_count              = var.max_use_count
  security_group_ids         = var.security_group_ids
  privileged                 = var.privileged
  region                     = var.region
  runner_repository          = var.runner_repository
  runner_version             = local.runner_version
  aws_plugin_version         = var.aws_plugin_version
  instance_role_profile_name = var.instance_role_profile_name
  install_cloudwatch_agent   = var.install_cloudwatch_agent
  cloudwatch_agent_json      = var.cloudwatch_agent_json
  enable_metrics_export      = var.enable_metrics_export
  metrics_export_endpoint    = var.metrics_export_endpoint
  default_docker_image       = var.default_docker_image
  runners_docker_section     = var.runners_docker_section
  runner_manager_ami         = var.runner_manager_ami
  usage_logger               = var.usage_logger
  acceptable_durations       = var.acceptable_durations
  node_exporter              = var.node_exporter
  create_key_pair            = var.create_key_pair

  name                        = var.metadata.name
  labels                      = module.labels.merged
  associate_public_ip_address = var.associate_public_ip_address
  instance_type               = var.instance_type
  encrypted                   = var.encrypted
  kms_key_id                  = var.kms_key_id
  volume_size                 = var.volume_size
  volume_type                 = var.volume_type
  throughput                  = var.throughput

  runner_wrapper = var.runner_wrapper
}
