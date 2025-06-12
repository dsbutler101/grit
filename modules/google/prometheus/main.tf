#######################
# METADATA VALIDATION #
#######################

module "validate_name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate_support" {
  source   = "../../internal/validation/support"
  use_case = "prometheus"
  use_case_support = tomap({
    "prometheus" = "experimental"
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

##########################
# PROMETHEUS PROD CONFIG #
##########################

locals {
  prometheus_server_tag = "prometheus-server"

  prometheus_image = "prom/prometheus:${var.prometheus_version}"

  data_device_id       = "persistent-data"
  persistent_data_path = "/mnt/disks/data"
  prometheus_volume    = "${local.persistent_data_path}/prometheus"
}

module "node_exporter" {
  source    = "../../internal/containers/node_exporter"
  image_tag = var.node_exporter_version
  port      = var.node_exporter_port
}

module "systemd_containers" {
  source = "../../internal/containers/systemd"

  containers = [
    {
      name  = "prometheus",
      image = local.prometheus_image
      ports = [
        "127.0.0.1:9090:9090"
      ]
      volumes = [
        "/etc/prometheus:/etc/prometheus",
        "${local.prometheus_volume}:/prometheus"
      ]
      command = <<EOF
      --web.enable-admin-api \
      --config.file=/etc/prometheus/config.yml \
      --storage.tsdb.path=/prometheus \
      --storage.tsdb.retention=7d
      EOF
      service_options = [{
        TimeoutStopSec = 30
      }]
    },
    module.node_exporter.container_config
  ]
}

data "cloudinit_config" "config" {
  gzip          = false
  base64_encode = false

  part {
    filename     = "cloud-config.yaml"
    content_type = "text/cloud-config"

    content = yamlencode({
      write_files = flatten([
        {
          path        = "/etc/prometheus/config.yml"
          owner       = "root:root"
          permissions = "0644"
          content     = local.prometheus_config_yml
        },
        {
          path        = "/etc/scripts/mount-prometheus-data-disk.sh"
          owner       = "root:root"
          permissions = "0700"
          content = templatefile("${path.module}/templates/mount-prometheus-data-disk.sh", {
            prometheus_image  = local.prometheus_image
            device_path       = "/dev/disk/by-id/google-${local.data_device_id}"
            mount_path        = local.persistent_data_path
            prometheus_volume = local.prometheus_volume
          })
        },
        module.systemd_containers.write_files
      ])

      runcmd = flatten([
        "/etc/scripts/mount-prometheus-data-disk.sh",
        module.systemd_containers.run_command
      ])
    })
  }
}

// This enforces Prometheus instance recreation when cloud_init configuration
// is changed
resource "terraform_data" "prometheus_server_replacement" {
  input = data.cloudinit_config.config.id
}

resource "google_compute_disk" "prometheus_data" {
  name   = "${var.metadata.name}-prometheus-data"
  labels = module.labels.merged

  zone = var.google_zone

  // Google Cloud API expects that the disk size will be a numerical string
  size = tostring(var.data_disk.size_gb)
  type = var.data_disk.disk_type
}

resource "google_compute_instance" "prometheus_server" {
  name         = "${var.metadata.name}-prometheus-${terraform_data.prometheus_server_replacement.output}"
  machine_type = var.machine_type

  lifecycle {
    // Because we use the attached volume, we must first remove the old instance
    // before creating the new one
    create_before_destroy = false

    // We don't need to worry about termination in case of Prometheus server
    // (as we do with the Runner Manager), so we can trigger replacement
    // on every configuration change.
    replace_triggered_by = [
      terraform_data.prometheus_server_replacement
    ]
  }

  metadata = {
    user-data              = data.cloudinit_config.config.rendered
    cos-update-strategy    = "update_disabled"
    enable-oslogin         = true
    block-project-ssh-keys = true
  }

  labels = merge(module.labels.merged, {
    purpose : local.prometheus_server_tag,
  })

  zone = var.google_zone

  tags = [
    local.prometheus_server_tag
  ]

  boot_disk {
    initialize_params {
      // Google Cloud API expects that the disk size will be a numerical string
      size = tostring(var.boot_disk.size_gb)
      type = var.boot_disk.disk_type

      image = "projects/cos-cloud/global/images/family/cos-stable"
    }
  }

  attached_disk {
    source      = google_compute_disk.prometheus_data.self_link
    device_name = local.data_device_id
  }

  network_interface {
    network    = var.vpc.id
    subnetwork = var.vpc.subnetwork_ids[var.prometheus_subnet_name]
    access_config {
      nat_ip = ""
    }
  }

  service_account {
    email = var.service_account_email
    scopes = [
      # Needed to allow Prometheus' gce_sd_config to discover
      # compute instances for monitoring
      "https://www.googleapis.com/auth/compute.readonly",

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

