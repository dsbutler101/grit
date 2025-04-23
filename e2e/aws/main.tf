terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.61.0"
    }
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = ">= 17.0.0"
    }
  }

  backend "http" {}
}

locals {
  metadata = {
    name = var.name
    labels = tomap({
      gitlab_project_id = var.gitlab_project_id
      env               = "grit-e2e"
    })
    min_support = "experimental"
  }
}

provider "gitlab" {}

module "gitlab" {
  source             = "../../modules/gitlab"
  metadata           = local.metadata
  url                = "https://gitlab.com"
  project_id         = var.gitlab_project_id
  runner_description = var.name
  runner_tags        = [var.runner_tag]
}

provider "aws" {}

// detect current configured region
data "aws_region" "current" {}

module "iam" {
  source   = "../../modules/aws/iam"
  metadata = local.metadata
}

module "vpc" {
  source   = "../../modules/aws/vpc"
  metadata = local.metadata

  // assumes every region we support has a second zone
  zone        = "${data.aws_region.current.name}b"
  cidr        = "10.0.0.0/16"
  subnet_cidr = "10.0.0.0/24"
}

module "security_groups" {
  source   = "../../modules/aws/security_groups"
  metadata = local.metadata

  vpc = local.vpc
}

module "fleeting" {
  source   = "../../modules/aws/fleeting"
  metadata = local.metadata

  vpc = local.vpc

  service              = "ec2"
  os                   = "linux"
  ephemeral_runner_ami = module.ami_lookup.ami_id
  instance_type        = var.ami_arch == "arm64" ? "t4g.medium" : "t2.medium"
  scale_min            = 1
  scale_max            = 10

  security_group_ids = [
    module.security_groups.fleeting_id,
  ]
}

module "ami_lookup" {
  source   = "../../modules/aws/ami_lookup"
  region   = data.aws_region.current.name
  metadata = local.metadata
  arch     = var.ami_arch
  os       = "linux"
  role     = "ephemeral"
}

module "cache" {
  source                = "../../modules/aws/cache"
  metadata              = local.metadata
  cache_object_lifetime = 2
}

module "runner" {
  source   = "../../modules/aws/runner"
  metadata = local.metadata
  vpc      = local.vpc
  iam      = local.iam
  fleeting = local.fleeting
  gitlab   = local.gitlab
  cache    = local.cache

  service               = "ec2"
  executor              = "docker-autoscaler"
  scale_min             = 1
  scale_max             = 10
  idle_percentage       = 10
  capacity_per_instance = 1

  security_group_ids = [
    module.security_groups.runner_manager_id,
  ]

  runner_version = "${var.runner_version}-1"
  runner_wrapper = {
    enabled = var.enable_runner_wrapper
  }
}

locals {
  iam = {
    enabled                    = module.iam.enabled
    fleeting_access_key_id     = module.iam.fleeting_access_key_id
    fleeting_secret_access_key = module.iam.fleeting_secret_access_key
  }

  vpc = {
    enabled    = module.vpc.enabled
    id         = module.vpc.id
    subnet_ids = module.vpc.subnet_ids
  }

  fleeting = {
    enabled                = module.fleeting.enabled
    ssh_key_pem_name       = module.fleeting.ssh_key_pem_name
    ssh_key_pem            = module.fleeting.ssh_key_pem
    autoscaling_group_name = module.fleeting.autoscaling_group_name
  }

  gitlab = {
    enabled      = module.gitlab.enabled
    url          = module.gitlab.url
    runner_token = module.gitlab.runner_token
  }

  cache = {
    enabled           = module.cache.enabled
    server_address    = module.cache.server_address
    bucket_name       = module.cache.bucket_name
    bucket_location   = module.cache.bucket_location
    access_key_id     = module.cache.access_key_id
    secret_access_key = module.cache.secret_access_key
  }
}
