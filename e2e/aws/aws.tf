variable "runner_token" {}
variable "name" {}
variable "job_id" {}

output "autoscaling_group_name" {
  value = module.fleeting.autoscaling_group_name
}

locals {
  metadata = {
    name = var.name
    labels = tomap({
      job_id = var.job_id
      env    = "grit-e2e"
    })
    min_support = "experimental"
  }
}

module "iam" {
  source   = "../../modules/aws/iam"
  metadata = local.metadata
}

module "vpc" {
  source   = "../../modules/aws/vpc/prod"
  metadata = local.metadata

  zone        = "us-east-1b"
  cidr        = "10.0.0.0/16"
  subnet_cidr = "10.0.0.0/24"
}

module "security_groups" {
  source   = "../../modules/aws/security_groups"
  metadata = local.metadata

  vpc_id = module.vpc.id
}

module "fleeting" {
  source   = "../../modules/aws/fleeting"
  metadata = local.metadata
  vpc      = local.vpc

  service       = "ec2"
  os            = "linux"
  ami           = module.ami_lookup.ami_id
  instance_type = "t2.medium"
  scale_min     = 1
  scale_max     = 10

  security_group_ids = [
    module.security_groups.fleeting.id,
  ]
}

module "ami_lookup" {
  source   = "../../modules/aws/ami_lookup"
  use_case = "aws-linux-ephemeral"
  region   = "us-east-1"
  metadata = local.metadata
}

module "s3_cache" {
  source                = "../../modules/aws/cache"
  metadata              = local.metadata
  cache_object_lifetime = 2
}

module "runner" {
  source   = "../../modules/aws/runner/prod"
  metadata = local.metadata
  vpc      = local.vpc
  iam      = local.iam
  fleeting = local.fleeting
  gitlab   = local.gitlab
  s3_cache = local.s3_cache

  service               = "ec2"
  executor              = "docker-autoscaler"
  scale_min             = 1
  scale_max             = 10
  idle_percentage       = 10
  capacity_per_instance = 1

  security_group_ids = [
    module.security_groups.runner_manager.id,
  ]
}

locals {
  iam = {
    fleeting_access_key_id     = module.iam.fleeting_access_key_id
    fleeting_secret_access_key = module.iam.fleeting_secret_access_key
  }

  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnet_id
  }

  fleeting = {
    ssh_key_pem_name       = module.fleeting.ssh_key_pem_name
    ssh_key_pem            = module.fleeting.ssh_key_pem
    autoscaling_group_name = module.fleeting.autoscaling_group_name
  }

  gitlab = {
    runner_token = var.runner_token
    url          = "https://gitlab.com"
  }

  s3_cache = {
    enabled           = module.s3_cache.enabled
    server_address    = module.s3_cache.server_address
    bucket_name       = module.s3_cache.bucket_name
    bucket_location   = module.s3_cache.bucket_location
    access_key_id     = module.s3_cache.access_key_id
    secret_access_key = module.s3_cache.secret_access_key
  }
}

