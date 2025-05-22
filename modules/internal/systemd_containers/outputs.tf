output "services" {
  description = "systemd service files for each container"
  value       = local.services
}

output "run_command" {
  description = "systemctl command to enable and start all containers"
  value       = local.run_command
}
