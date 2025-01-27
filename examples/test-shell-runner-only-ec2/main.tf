variable "runner_token" {
  default = "MY_RUNNER_TOKEN"
}

locals {
  metadata = {
    name        = "littlerunner"
    min_support = "none"
    labels = {
      "env" = "test"
    }
  }
}

module "runner" {
  source = "../../modules/aws/runner/test"

  metadata = local.metadata

  service = "ec2"
  gitlab = {
    runner_token = var.runner_token
  }
  vpc = {
    id        = "vpc-0d119da238d878eef"
    subnet_id = "subnet-0bd3ab8c221e14bfc"
  }

  security_group_ids = [module.security_groups.runner_manager.id]
}

module "security_groups" {
  source = "../../modules/aws/security_groups"

  metadata = local.metadata

  vpc_id = "vpc-0d119da238d878eef"
}
