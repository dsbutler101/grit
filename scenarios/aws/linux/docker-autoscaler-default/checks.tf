module "check_concurrent_and_capacity" {
  source  = "../../../../modules/internal/validation/fail_validation"
  message = var.max_instances >= var.concurrent / var.capacity_per_instance ? "" : "max_instances needs to be higher to fit the defined number of concurrent with the given capacity_per_instance"
}
