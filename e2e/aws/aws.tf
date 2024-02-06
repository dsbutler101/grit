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
  source   = "../../modules/aws/iam/prod"
  metadata = local.metadata
}

module "vpc" {
  source   = "../../modules/aws/vpc/prod"
  metadata = local.metadata

  zone        = "us-east-1b"
  cidr        = "10.0.0.0/16"
  subnet_cidr = "10.0.0.0/24"
}

module "fleeting" {
  source   = "../../modules/aws/fleeting/prod"
  metadata = local.metadata
  vpc      = local.vpc

  service       = "ec2"
  os            = "linux"
  ami           = "ami-0735db9b38fcbdb39"
  instance_type = "t2.medium"
  scale_min     = 1
  scale_max     = 10
}

module "runner" {
  source   = "../../modules/aws/runner/prod"
  metadata = local.metadata
  vpc      = local.vpc
  iam      = local.iam
  fleeting = local.fleeting
  gitlab   = local.gitlab

  service               = "ec2"
  executor              = "docker-autoscaler"
  scale_min             = 1
  scale_max             = 10
  idle_percentage       = 10
  capacity_per_instance = 1
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
}

