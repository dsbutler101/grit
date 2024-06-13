module "vpc" {
  source = "../../../../modules/google/vpc/prod"

  metadata = local.metadata

  google_region = var.google_region

  subnetworks = {
    runner-manager    = "10.0.0.0/29"
    ephemeral-runners = "10.1.0.0/21"
  }
}

module "iam" {
  source = "../../../../modules/google/iam/prod"

  metadata = local.metadata
}

module "cache" {
  source = "../../../../modules/google/cache/prod"

  metadata = local.metadata

  bucket_location = var.google_region

  service_account_emails = [
    module.iam.service_account_email
  ]
}

module "fleeting" {
  source = "../../../../modules/google/fleeting/prod"

  metadata = local.metadata

  google_project = var.google_project
  google_zone    = var.google_zone

  fleeting_service = "gce"

  service_account_email = module.iam.service_account_email

  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnetwork_ids["ephemeral-runners"]
  }

  manager_subnet_cidr = module.vpc.subnetwork_cidrs["runner-manager"]

  disk_type    = var.ephemeral_runner.disk_type
  disk_size_gb = var.ephemeral_runner.disk_size
  machine_type = var.ephemeral_runner.machine_type
  source_image = var.ephemeral_runner.source_image
}

module "runner" {
  source = "../../../../modules/google/runner/prod"

  metadata = local.metadata

  google_project = var.google_project
  google_zone    = var.google_zone

  service_account_email = module.iam.service_account_email

  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnetwork_ids["runner-manager"]
  }

  gitlab_url   = var.gitlab_url
  runner_token = var.runner_token

  executor = "docker-autoscaler"

  cache_gcs_bucket = module.cache.bucket_name

  fleeting_instance_group_name = module.fleeting.instance_group_name

  runners_global_section = var.runners_global_section
  runners_docker_section = var.runners_docker_section

  machine_type = var.runner_machine_type

  concurrent     = var.concurrent
  check_interval = 3

  request_concurrency = var.concurrent > 10 ? 10 : var.concurrent

  capacity_per_instance = var.capacity_per_instance
  max_instances         = var.max_instances
  max_use_count         = var.max_use_count
  autoscaling_policies  = concat([local.required_autoscaling_policy], var.autoscaling_policies)
}
