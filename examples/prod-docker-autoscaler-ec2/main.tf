# export GITLAB_TOKEN=<YOUR_GITLAB_TOKEN>

terraform {
  required_providers {
    gitlab = {
      source = "gitlabhq/gitlab"
    }
  }
}

locals {
  # Metadata is common input to all modules.
  metadata = {
    name = "my-autoscaling-runner"
    labels = tomap({
      env = "grit-e2e"
    })
    min_support = "experimental"
  }
}

# The IAM module creates a service account for runner to access
# ephemeral VMs.
module "iam" {
  source   = "../../modules/aws/iam/prod"
  metadata = local.metadata
}

# The VPC module creates an independent VPC to isolate this runner
# installation from others.
module "vpc" {
  source   = "../../modules/aws/vpc/prod"
  metadata = local.metadata

  zone        = "us-east-1b"
  cidr        = "10.0.0.0/16"
  subnet_cidr = "10.0.0.0/24"
}

# The fleeting module creates an instance group and SSH key to access
# it. This is used as job running capacity for the autoscaling runner.
module "fleeting" {
  source   = "../../modules/aws/fleeting/prod"
  metadata = local.metadata

  # The outputs of the VPC module are passed into fleeting so it knows
  # where to put the VMs.
  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnet_id
  }

  service       = "ec2"
  os            = "linux"
  ami           = "ami-0735db9b38fcbdb39"
  instance_type = "t2.medium"
  scale_min     = 1
  scale_max     = 10
}

# The gitlab modules will register the created runner to GitLab as a
# project runner.
module "gitlab" {
  source   = "../../modules/gitlab/prod"
  metadata = local.metadata

  url                = "https://gitlab.com"
  project_id         = "39258790"
  runner_description = "An example docker-autoscaler runner on EC2"
  runner_tags        = []
}

module "runner" {
  source   = "../../modules/aws/runner/prod"
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
}
