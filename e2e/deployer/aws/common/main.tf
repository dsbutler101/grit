// detect current configured region
data "aws_region" "current" {}

locals {
  metadata = {
    name        = var.name
    min_support = "none"
    labels = {
      "env" = "grit-e2e"
    }
  }
}

module "gitlab" {
  source             = "../../../../modules/gitlab"
  metadata           = local.metadata
  url                = "https://gitlab.com"
  project_id         = var.gitlab_project_id
  runner_description = var.name
  runner_tags        = var.runner_tag_list
}

module "vpc" {
  source   = "../../../../modules/aws/vpc"
  metadata = local.metadata

  // assumes every region we support has a second zone
  zone        = "${data.aws_region.current.name}b"
  cidr        = "10.0.0.0/16"
  subnet_cidr = "10.0.0.0/24"
}

module "security_groups" {
  source   = "../../../../modules/aws/security_groups"
  metadata = local.metadata

  vpc = local.vpc
}

module "runner" {
  source   = "../../../../modules/aws/runner"
  metadata = local.metadata
  vpc      = local.vpc
  gitlab   = local.gitlab

  service               = "ec2"
  executor              = "shell"
  capacity_per_instance = 2

  # Needed because https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit/-/blob/b1a60990e4e6b1b5a4a5c7ff588fd88ddd8c1428/modules/aws/runner/internal/ec2/ec2.tf#L26 multiplies.
  scale_max = 1

  security_group_ids = [
    module.security_groups.runner_manager_id,
  ]

  instance_type = "t3.small"

  runner_version = "${var.runner_version}-1"
  runner_wrapper = {
    enabled     = true
    socket_path = "tcp://localhost:7777"
  }
  create_key_pair = {}
}

locals {
  vpc = {
    enabled    = true
    id         = module.vpc.id
    subnet_ids = module.vpc.subnet_ids
  }

  gitlab = {
    enabled      = true
    url          = module.gitlab.url
    runner_token = module.gitlab.runner_token
  }
}
