module "vpc" {
  source = "../../../../modules/google/vpc/prod"

  metadata = local.metadata

  google_region = var.google_region

  subnetworks = local.subnetworks
}

module "iam" {
  source = "../../../../modules/google/iam/prod"

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

module "prometheus-iam" {
  source = "../../../../modules/google/iam/prod"

  metadata = merge(local.metadata, {
    name = "${local.metadata.name}-prom"
  })
}

module "prometheus" {
  count = var.prometheus.enabled ? 1 : 0

  source = "../../../../modules/google/prometheus/prod"

  metadata = local.metadata

  google_project = var.google_project
  google_zone    = var.google_zone

  service_account_email = module.prometheus-iam.service_account_email

  vpc = {
    id        = module.vpc.id
    subnet_id = module.vpc.subnetwork_ids["prometheus"]
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
