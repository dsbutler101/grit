locals {
  check_scale_min_fail_message             = "scale_min is required for the autoscaling instance and docker-autoscaler executors"
  check_scale_max_fail_message             = "scale_max is required for the autoscaling instance and docker-autoscaler executors"
  check_idle_percentage_fail_message       = "idle_percentage is required for the autoscaling instance and docker-autoscaler executors"
  check_capacity_per_instance_fail_message = "capacity_per_instance is required for the autoscaling instance and docker-autoscaler executors"
}

module "check_scale_min" {
  source  = "../fail_validation"
  message = var.scale_min == -1 && (var.executor == "instance" || var.executor == "docker-autoscaler") ? local.check_scale_min_fail_message : ""
}

module "check_scale_max" {
  source  = "../fail_validation"
  message = var.scale_max == -1 && (var.executor == "instance" || var.executor == "docker-autoscaler") ? local.check_scale_max_fail_message : ""
}

module "check_idle_percentage" {
  source  = "../fail_validation"
  message = var.idle_percentage == -1 && (var.executor == "instance" || var.executor == "docker-autoscaler") ? local.check_idle_percentage_fail_message : ""
}

module "check_capacity_per_instance" {
  source  = "../fail_validation"
  message = var.capacity_per_instance == -1 && (var.executor == "instance" || var.executor == "docker-autoscaler") ? local.check_capacity_per_instance_fail_message : ""
}
