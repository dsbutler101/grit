output "grit_runner_managers" {
  value = tomap({
    "runner-manager" = tomap({
      instance_name   = module.runner.instance_name
      instance_id     = module.runner.instance_id
      address         = "${module.runner.external_ip}:22"
      wrapper_address = module.runner.wrapper_address
    })
  })
  sensitive = true
}
