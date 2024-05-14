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
    periods            = var.autoscaling_policies.periods
    timezone           = var.autoscaling_policies.timezone
    scale_min          = var.autoscaling_policies.scale_min
    idle_time          = var.autoscaling_policies.idle_time
    scale_factor       = var.autoscaling_policies.scale_factor
    scale_factor_limit = var.autoscaling_policies.scale_factor_limit
  }
}
