locals {
  default_labels = {
    managed = "grit"
  }

  metadata = {
    name        = var.name
    labels      = merge(var.labels, local.default_labels)
    min_support = "experimental"
  }

  required_autoscaling_policy = {
    periods            = ["* * * * *"]
    timezone           = ""
    scale_min          = 10
    idle_time          = "20m0s"
    scale_factor       = 0.2
    scale_factor_limit = 100
  }

  runner_metrics_port = 9252
  node_expoter_port   = 9100

  subnetworks_base = {
    runner-manager    = "10.0.0.0/29"
    ephemeral-runners = "10.1.0.0/21"
  }

  subnetworks = merge(
    local.subnetworks_base,
    var.prometheus.enabled ? { prometheus : "10.0.0.8/29" } : {},
  )
}
