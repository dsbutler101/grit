###############################
# INTERNAL EC2 INSTANCE GROUP #
###############################

module "common" {
  source = "./common"

  license_arn                      = try(module.macos[0].license-config-arn, "")
  jobs-host-resource-group-outputs = try(module.macos[0].jobs-host-resource-group-outputs, {})

  scale_min = var.scale_min
  scale_max = var.scale_max

  storage_type               = var.storage_type
  storage_size               = var.storage_size
  storage_throughput         = var.storage_throughput
  ami_id                     = var.ami
  instance_type              = var.instance_type
  vpc_id                     = var.vpc.id
  subnet_id                  = var.vpc.subnet_id
  security_group_ids         = var.security_group_ids
  install_cloudwatch_agent   = var.install_cloudwatch_agent
  cloudwatch_agent_json      = var.cloudwatch_agent_json
  instance_role_profile_name = var.instance_role_profile_name
  mixed_instances_policy     = var.mixed_instances_policy
  ebs_encryption             = var.ebs_encryption
  kms_key_arn                = var.kms_key_arn

  labels = var.labels
  name   = var.name
}

module "macos" {
  count  = var.os == "macos" ? 1 : 0
  source = "./macos"

  license_count_per_asg = var.macos_license_count_per_asg
  cores_per_license     = var.macos_cores_per_license

  labels = var.labels
  name   = var.name
}
