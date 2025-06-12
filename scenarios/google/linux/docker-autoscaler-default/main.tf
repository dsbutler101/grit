module "vpc" {
  source = "../../../../modules/google/vpc"

  metadata = local.metadata

  google_region = var.google_region

  subnetworks = local.subnetworks
}

module "iam" {
  source = "../../../../modules/google/iam"

  metadata = local.metadata
}

module "cache" {
  source = "../../../../modules/google/cache"

  metadata = local.metadata

  bucket_location = var.google_region

  service_account_emails = [
    module.iam.service_account_email
  ]
}

module "fleeting" {
  source = "../../../../modules/google/fleeting"

  metadata = local.metadata

  google_project     = var.google_project
  google_zone        = var.google_zone
  subnetwork_project = var.google_project

  fleeting_service = "gce"

  service_account_email = module.iam.service_account_email

  vpc = {
    enabled          = module.vpc.enabled
    id               = module.vpc.id
    subnetwork_ids   = module.vpc.subnetwork_ids
    subnetwork_cidrs = module.vpc.subnetwork_cidrs
  }

  disk_type    = var.ephemeral_runner.disk_type
  disk_size_gb = var.ephemeral_runner.disk_size
  machine_type = var.ephemeral_runner.machine_type
  source_image = var.ephemeral_runner.source_image
}

module "runner_version_lookup" {
  source = "../../../../modules/gitlab/runner_version_lookup"

  metadata = local.metadata

  skew = var.runner_version_skew
}

module "runner" {
  source = "../../../../modules/google/runner"

  metadata = local.metadata

  google_project     = var.google_project
  google_zone        = var.google_zone
  subnetwork_project = var.google_project

  service_account_email = module.iam.service_account_email

  vpc = {
    enabled          = module.vpc.enabled
    id               = module.vpc.id
    subnetwork_ids   = module.vpc.subnetwork_ids
    subnetwork_cidrs = module.vpc.subnetwork_cidrs
  }

  node_exporter = {
    port : local.node_expoter_port
  }

  listen_address = "" // This is set to disable the deprecated listen_address and use runner_metrics_listener instead
  runner_metrics_listener = {
    address = "0.0.0.0"
    port    = local.runner_metrics_port
  }

  gitlab_url   = var.gitlab_url
  runner_token = var.runner_token

  runner_version_lookup = {
    skew           = module.runner_version_lookup.skew
    runner_version = module.runner_version_lookup.runner_version
  }

  executor = "docker-autoscaler"

  cache_gcs_bucket = module.cache.bucket_name

  fleeting_instance_group_name = module.fleeting.instance_group_name

  runners_global_section = var.runners_global_section
  runners_docker_section = var.runners_docker_section

  machine_type = var.runner_machine_type
  disk_type    = var.runner_disk_type

  concurrent     = var.concurrent
  check_interval = 3

  request_concurrency = var.concurrent > 10 ? 10 : var.concurrent

  capacity_per_instance = var.capacity_per_instance
  max_instances         = var.max_instances
  max_use_count         = var.max_use_count
  autoscaling_policies  = concat([local.required_autoscaling_policy], var.autoscaling_policies)

  runner_manager_additional_firewall_rules = !var.prometheus.enabled ? {} : {
    prometheus = {
      direction = "INGRESS"
      priority  = 1000
      allow = [
        {
          protocol = "tcp"
          ports = [
            local.runner_metrics_port,
            local.node_expoter_port,
          ]
        }
      ]
      source_ranges = [
        module.vpc.subnetwork_cidrs["prometheus"]
      ]
    }
  }
}

module "prometheus_iam" {
  source = "../../../../modules/google/iam"

  metadata = merge(local.metadata, {
    name = "${local.metadata.name}-p"
  })
}

module "prometheus" {
  count = var.prometheus.enabled ? 1 : 0

  source = "../../../../modules/google/prometheus"

  metadata = local.metadata

  google_project = var.google_project
  google_zone    = var.google_zone

  service_account_email = module.prometheus_iam.service_account_email

  vpc = {
    enabled          = module.vpc.enabled
    id               = module.vpc.id
    subnetwork_ids   = module.vpc.subnetwork_ids
    subnetwork_cidrs = module.vpc.subnetwork_cidrs
  }

  node_exporter_port = local.node_expoter_port

  mimir = var.prometheus.mimir

  prometheus_external_labels = var.prometheus.external_labels

  runner_manager_nodes = {
    filter = "labels.purpose = \"gitlab-runner-manager\""

    exporter_ports = {
      runner_manager = local.runner_metrics_port
      node_exporter  = local.node_expoter_port
    }

    custom_relabel_configs = var.prometheus.custom_relabel_configs

    instance_labels_to_include = var.prometheus.instance_labels_to_include
  }
}
