module "check-concurrent" {
  source  = "../../../../modules/internal/validation/fail_validation"
  message = var.concurrent <= 2000 ? "" : "Configuration (especially network size) will not be able to handle more than 2000 concurrent jobs"
}

module "check-max-instances" {
  source  = "../../../../modules/internal/validation/fail_validation"
  message = var.autoscaling_policy.scale_max <= 1000 ? "" : "Fleeting plugin for AWS will not allow to manage more than 1000 instances at once"
}

 module "check-concurrent-and-capacity" {
   source  = "../../../../modules/internal/validation/fail_validation"
   message = var.autoscaling_policy.scale_max >= var.concurrent / var.capacity_per_instance ? "" : "max_instances needs to be higher to fit the defined number of concurrent with the given capacity_per_instance"
 }
