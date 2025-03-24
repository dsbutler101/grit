module "vpc" {
  source   = "../../../../modules/aws/vpc"
  metadata = local.metadata

  zone = var.aws_zone

  cidr        = "10.0.0.0/16"
  subnet_cidr = "10.0.0.0/24"
}

module "iam" {
  source = "../../../../modules/aws/iam"

  metadata = local.metadata
}

module "ami_lookup" {
  source   = "../../../../modules/aws/ami_lookup"
  use_case = "aws-linux-ephemeral"
  region   = var.aws_region
  metadata = local.metadata
}

module "fleeting" {
  source = "../../../../modules/aws/fleeting"

  metadata = local.metadata

  service = "ec2"
  os      = "linux"

  vpc = {
    id         = module.vpc.id
    subnet_ids = module.vpc.subnet_ids
  }

  security_group_ids = [module.security_groups.fleeting.id]

  instance_type        = var.ephemeral_runner.machine_type
  ephemeral_runner_ami = var.ephemeral_runner.source_image != "" ? var.ephemeral_runner.source_image : module.ami_lookup.ami_id
  scale_min            = var.autoscaling_policy.scale_min
  scale_max            = var.max_instances
}

module "cache" {
  source = "../../../../modules/aws/cache"

  metadata = local.metadata
}

module "runner" {
  source = "../../../../modules/aws/runner"

  metadata = local.metadata

  vpc = {
    id         = module.vpc.id
    subnet_ids = module.vpc.subnet_ids
  }
  iam = {
    fleeting_access_key_id     = module.iam.fleeting_access_key_id
    fleeting_secret_access_key = module.iam.fleeting_secret_access_key
  }
  fleeting = {
    ssh_key_pem_name       = module.fleeting.ssh_key_pem_name
    ssh_key_pem            = module.fleeting.ssh_key_pem
    autoscaling_group_name = module.fleeting.autoscaling_group_name
  }
  gitlab = {
    runner_token = module.gitlab.runner_token
    url          = module.gitlab.url
  }
  s3_cache = {
    enabled           = module.cache.enabled
    server_address    = module.cache.server_address
    bucket_name       = module.cache.bucket_name
    bucket_location   = module.cache.bucket_location
    access_key_id     = module.cache.access_key_id
    secret_access_key = module.cache.secret_access_key
  }

  service               = "ec2"
  executor              = "docker-autoscaler"
  scale_min             = var.autoscaling_policy.scale_min
  scale_max             = var.max_instances
  idle_percentage       = var.autoscaling_policy.scale_factor
  capacity_per_instance = var.capacity_per_instance
  max_use_count         = var.max_use_count
  region                = var.aws_region

  security_group_ids = [module.security_groups.runner_manager.id]
}

module "security_groups" {
  source   = "../../../../modules/aws/security_groups"
  metadata = local.metadata

  vpc_id = module.vpc.id
}

module "gitlab" {
  source   = "../../../../modules/gitlab"
  metadata = local.metadata

  url                = var.gitlab_url
  project_id         = var.gitlab_project_id
  runner_description = var.runner_description
  runner_tags        = var.runner_tags
}
