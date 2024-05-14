# TODO: create variables and move them to main.tf

module "vpc" {
  source   = "../../../../modules/aws/vpc/prod"
  metadata = local.metadata

  zone        = var.aws_zone

  cidr        = "10.0.0.0/16"
  subnet_cidr = "10.0.0.0/24"
}

module "iam" {
  source   = "../../../../modules/aws/iam/prod"

  metadata = local.metadata
}

module "fleeting" {
  source   = "../../../../modules/aws/fleeting/prod"

  metadata = local.metadata

  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnet_id
  }

  service       = "ec2"
  os            = "linux"
  ami           = var.ephemeral_runner.source_image
  instance_type = var.ephemeral_runner.machine_type
  scale_min     = 1
  scale_max     = 10

  security_group_ids = [module.security_groups.fleeting.id]
}

# The gitlab modules will register the created runner to GitLab as a
# project runner.
module "gitlab" {
  source   = "../.local/grit/modules/gitlab/prod"
  metadata = local.metadata

  url                = "https://gitlab.com"
  project_id         = "56778975"
  runner_description = "An example docker-autoscaler runner on EC2"
  runner_tags        = ["my-runner"]
}

module "runner" {
  source   = "../.local/grit/modules/aws/runner/prod"
  metadata = local.metadata

  # All the module outputs are ultimately fed into the runner module
  # to create and configure the runner manager.
  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnet_id
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

  service               = "ec2"
  executor              = "docker-autoscaler"
  scale_min             = 1
  scale_max             = 10
  idle_percentage       = 10
  capacity_per_instance = 1

  security_group_ids = [module.security_groups.runner_manager.id]
}

module "security_groups" {
  source   = "../.local/grit/modules/aws/security_groups/prod"
  metadata = local.metadata

  vpc_id = module.vpc.id

}
