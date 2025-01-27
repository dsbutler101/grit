locals {
  runner_manager_nodes_gce_sd_config_base = {
    project = var.google_project
    zone    = var.google_zone
    filter  = var.runner_manager_nodes.filter
  }

  runner_manager_nodes_relabel_configs_base = [
    {
      target_label = "instance"
      source_labels = [
        "__meta_gce_instance_name"
      ]
    },
    {
      target_label = "zone"
      source_labels = [
        "__meta_gce_instance_zone"
      ]
      // For the zone we're interested only in the zone identifier, not the full Google Cloud URI
      regex       = ".*/([^/]+)$"
      replacement = "$1"
    }
  ]

  runner_manager_nodes_relabel_configs = concat(
    local.runner_manager_nodes_relabel_configs_base,
    var.runner_manager_nodes.custom_relabel_configs,
    [
      for label in var.runner_manager_nodes.instance_labels_to_include :
      {
        target_label = label
        source_labels = [
          "__meta_gce_label_${label}"
        ]
      }
    ]
  )

  mimir_remote_write = var.mimir == null ? [] : (
    var.mimir.url == "" ? [] : [
      {
        url = var.mimir.url
        headers = {
          X-Scope-OrgID = var.mimir.tenant
        }
      }
    ]
  )

  prometheus_config_remote_write = concat(
    local.mimir_remote_write
  )

  prometheus_config = {
    global = {
      scrape_interval = "15s"
      external_labels = var.prometheus_external_labels
    }

    remote_write = local.prometheus_config_remote_write

    scrape_configs = [
      // Track self
      {
        job_name     = "prometheus"
        metrics_path = "/metrics"
        static_configs = [
          {
            targets = [
              "127.0.0.1:9090"
            ]
          }
        ]
      },

      // Node exporter
      {
        job_name     = "node"
        metrics_path = "/metrics"
        static_configs = [
          {
            // Points node exporter running on the host where Prometheus
            // is deployed
            targets = [
              "172.17.0.1:${var.node_exporter_port}"
            ]
          }
        ]
        gce_sd_configs = [
          merge(local.runner_manager_nodes_gce_sd_config_base, {
            port = var.runner_manager_nodes.exporter_ports.node_exporter
          })
        ]
        relabel_configs = local.runner_manager_nodes_relabel_configs
      },

      // Runner manager internal exporter
      {
        job_name     = "runners-manager"
        metrics_path = "/metrics"
        gce_sd_configs = [
          merge(local.runner_manager_nodes_gce_sd_config_base, {
            port = var.runner_manager_nodes.exporter_ports.runner_manager
          })
        ]
        relabel_configs = local.runner_manager_nodes_relabel_configs
      }
    ]
  }

  prometheus_config_yml = yamlencode(local.prometheus_config)
}
