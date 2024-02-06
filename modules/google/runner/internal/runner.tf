locals {
  runner_manager_tag = "gitlab-runner-manager"
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

            name       = var.name
            gitlab_url = var.gitlab_url

            runner_token   = google_kms_secret_ciphertext.runner-token.ciphertext
            runner_ssh_key = google_kms_secret_ciphertext.runner-ssh-key.ciphertext
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
            listen_address = var.listen_address
          })
        },
        {
          path        = "/etc/gitlab-runner/config-template.toml"
          owner       = "root:root"
          permissions = "0600"
          content = templatefile("${path.module}/templates/config-template.toml", {
            request_concurrency = var.request_concurrency

            cache_gcs_bucket = var.cache_gcs_bucket

            runners_global_section = var.runners_global_section
            runners_docker_section = var.runners_docker_section
          })
        },
        {
          path        = "/etc/systemd/system/gitlab-runner.service"
          owner       = "root:root"
          permissions = "0644"
          content = templatefile("${path.module}/templates/gitlab-runner.service", {
            gitlab_runner_image = "registry.gitlab.com/gitlab-org/gitlab-runner:alpine-${var.runner_version}"
          })
        },
      ]

      runcmd = [
        "systemctl daemon-reload",
        "systemctl start gitlab-runner.service",
      ]
    })
  }
}

resource "google_compute_instance" "runner-manager" {
  name         = "${var.name}-runner-manager"
  machine_type = var.machine_type

  metadata = {
    user-data      = data.cloudinit_config.config.rendered
    enable-oslogin = true
  }

  labels = var.labels

  tags = [
    local.runner_manager_tag
  ]

  boot_disk {
    initialize_params {
      type  = var.disk_type
      image = "projects/cos-cloud/global/images/family/cos-stable"
      size  = var.disk_size_gb
    }
  }

  network_interface {
    network    = var.vpc.id
    subnetwork = var.vpc.subnet_id
    access_config {
      nat_ip = ""
    }
  }

  service_account {
    email = var.service_account_email
    scopes = [
      # Needed for secrets decryption through Google KMS
      "https://www.googleapis.com/auth/cloudkms",

      # Needed for signing GCS URLs for cache
      "https://www.googleapis.com/auth/iam",

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
