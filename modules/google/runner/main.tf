#######################
# METADATA VALIDATION #
#######################

module "validate_name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate_support" {
  source   = "../../internal/validation/support"
  use_case = "runner"
  use_case_support = tomap({
    "runner" = "experimental"
  })
  min_support = var.metadata.min_support
}

##################
# DEFAULT LABELS #
##################

module "labels" {
  source = "../../internal/labels"

  name              = var.metadata.name
  additional_labels = var.metadata.labels
}

######################
# RUNNER PROD CONFIG #
######################

locals {
  runner_manager_name = "${var.metadata.name}-runner-manager"
  runner_manager_tag  = "gitlab-runner-manager"
  runner_manager_tags = concat([local.runner_manager_tag], var.additional_tags)

  use_autoscaling = var.executor == "docker-autoscaler" || var.executor == "instance"
  use_docker      = var.executor == "docker-autoscaler" || var.executor == "docker"

  autoscaling_policies = [
    for p in var.autoscaling_policies : {
      periods            = join(", ", formatlist("%q", p.periods))
      timezone           = p.timezone
      idle_count         = p.scale_min * var.capacity_per_instance
      idle_time          = p.idle_time
      scale_factor       = p.scale_factor
      scale_factor_limit = p.scale_factor_limit
    }
  ]

  runner_manager_machine_types_map = {
    0   = "c2-standard-4",
    300 = "c2-standard-8"
    600 = "c2-standard-16"
    900 = "c2-standard-30"
  }

  runner_manager_machine_type = var.concurrent < 300 ? local.runner_manager_machine_types_map[0] : (
    var.concurrent < 600 ? local.runner_manager_machine_types_map[300] : (
      var.concurrent < 900 ? local.runner_manager_machine_types_map[600] : local.runner_manager_machine_types_map[900]
    )
  )

  // These few lines are added to handle listen_address deprecation and backward compatibility
  //
  // DEPRECATED: we should switch to use runner_metrics_listener variable instead of listen_address
  metrics_listen_address_and_port = split(":", var.listen_address)
  metrics_listener_address        = var.listen_address != "" ? local.metrics_listen_address_and_port[0] : var.runner_metrics_listener.address
  metrics_listener_port           = var.listen_address != "" ? local.metrics_listen_address_and_port[1] : var.runner_metrics_listener.port
}

data "cloudinit_config" "config" {
  gzip          = false
  base64_encode = false

  part {
    filename     = "cloud-config.yaml"
    content_type = "text/cloud-config"

    content = yamlencode({
      write_files = [
        {
          path        = "/etc/gitlab-runner/entrypoint.sh"
          owner       = "root:root"
          permissions = "0755"
          content = templatefile("${path.module}/templates/entrypoint.sh", {
            kms_key = google_kms_crypto_key.default.id

            name       = var.metadata.name
            gitlab_url = var.gitlab_url

            runner_token   = google_kms_secret_ciphertext.runner_token.ciphertext
            runner_ssh_key = google_kms_secret_ciphertext.runner_ssh_key.ciphertext

            use_autoscaling                       = local.use_autoscaling
            fleeting_googlecompute_plugin_version = var.fleeting_googlecompute_plugin_version

            https_proxy = var.https_proxy
            http_proxy  = var.http_proxy
            no_proxy    = var.no_proxy
          })
        },
        {
          path        = "/etc/gitlab-runner/config.toml"
          owner       = "root:root"
          permissions = "0600"
          content = templatefile("${path.module}/templates/config.toml", {
            concurrent     = var.concurrent
            check_interval = var.check_interval
            log_level      = var.log_level
            log_format     = "text"
            listen_address = "${local.metrics_listener_address}:${local.metrics_listener_port}"
          })
        },
        {
          path        = "/etc/gitlab-runner/config-template.toml"
          owner       = "root:root"
          permissions = "0600"
          content = templatefile("${path.module}/templates/config-template.toml", {
            request_concurrency = var.request_concurrency

            cache_gcs_bucket = var.cache_gcs_bucket

            use_autoscaling = local.use_autoscaling
            use_docker      = local.use_docker

            executor = var.executor

            runners_global_section = var.runners_global_section
            runners_docker_section = var.runners_docker_section
            default_docker_image   = var.default_docker_image

            fleeting_google_project      = var.google_project
            fleeting_google_zone         = var.google_zone
            fleeting_instance_group_name = var.fleeting_instance_group_name

            capacity_per_instance = var.capacity_per_instance
            max_use_count         = var.max_use_count
            max_instances         = var.max_instances

            autoscaling_policies = local.autoscaling_policies
          })
        },
        {
          path        = "/etc/systemd/system/gitlab-runner.service"
          owner       = "root:root"
          permissions = "0644"
          content = templatefile("${path.module}/templates/gitlab-runner.service", {
            gitlab_runner_image = "${var.runner_registry}:alpine-${var.runner_version}"
            runner_metrics_port = local.metrics_listener_port
            additional_volumes  = var.additional_volumes
          })
        },
        {
          path        = "/etc/systemd/system/node-exporter.service"
          owner       = "root:root"
          permissions = "0644"
          content = templatefile("${path.module}/templates/node-exporter.service", {
            node_exporter_image = "prom/node-exporter:${var.node_exporter.version}"
            node_exporter_port  = var.node_exporter.port
          })
        },
      ]

      runcmd = [
        "systemctl daemon-reload",
        "systemctl enable node-exporter.service",
        "systemctl start node-exporter.service",
        "systemctl enable gitlab-runner.service",
        "systemctl start gitlab-runner.service",
      ]
    })
  }
}

resource "google_compute_instance" "runner_manager" {
  name         = local.runner_manager_name
  machine_type = var.machine_type != "" ? var.machine_type : local.runner_manager_machine_type

  metadata = {
    user-data           = data.cloudinit_config.config.rendered
    enable-oslogin      = true
    cos-update-strategy = "update_disabled"
  }

  labels = merge(module.labels.merged, {
    purpose = local.runner_manager_tag
  })

  zone = var.google_zone

  tags = local.runner_manager_tags

  boot_disk {
    initialize_params {
      type  = var.disk_type
      image = var.disk_image
      size  = var.disk_size_gb
    }
  }

  network_interface {
    network            = var.vpc.id
    subnetwork         = var.vpc.subnet_id
    subnetwork_project = var.subnetwork_project

    dynamic "access_config" {
      for_each = var.access_config_enabled ? [1] : []
      content {
        nat_ip = ""
      }
    }
  }

  service_account {
    email = var.service_account_email
    scopes = [
      # Needed for secrets decryption through Google KMS
      "https://www.googleapis.com/auth/cloudkms",

      # Needed for signing GCS URLs for cache
      "https://www.googleapis.com/auth/iam",

      # Needed for managing instances through the Instance Group Manager
      "https://www.googleapis.com/auth/compute",

      # The default scopes present if not defined explicitly as above
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring.write",
      "https://www.googleapis.com/auth/pubsub",
      "https://www.googleapis.com/auth/service.management.readonly",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/trace.append"
    ]
  }
}

resource "google_compute_address" "runner_manager" {
  name         = local.runner_manager_name
  address_type = var.address_type
  subnetwork   = var.address_type == "INTERNAL" ? var.vpc.subnet_id : null
}
