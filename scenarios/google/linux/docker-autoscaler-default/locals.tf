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
}
