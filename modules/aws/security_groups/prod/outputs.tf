output "fleeting" {
  description = "Security group for the ephemeral fleeting VMs"
  value       = module.security_groups.fleeting
}

output "runner_manager" {
  description = "Security group for the runner manager"
  value       = module.security_groups.runner_manager
}
