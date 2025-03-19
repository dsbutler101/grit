module "check_concurrent" {
  source  = "../../../../modules/internal/validation/fail_validation"
  message = var.concurrent <= 2000 ? "" : "Configuration (especially network size) will not be able to handle more than 2000 concurrent jobs"
}

module "check_max_instances" {
  source  = "../../../../modules/internal/validation/fail_validation"
  message = var.max_instances <= 1000 ? "" : "Fleeting plugin for Google will not allow to manage more than 1000 instances at once"
}

module "check_concurrent_and_capacity" {
  source  = "../../../../modules/internal/validation/fail_validation"
  message = var.max_instances >= var.concurrent / var.capacity_per_instance ? "" : "max_instances needs to be higher to fit the defined number of concurrent with the given capacity_per_instance"
}
