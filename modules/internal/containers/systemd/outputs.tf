output "write_files" {
  description = "systemd service files for each container for use with cloud-init's `write_files` module."
  value       = local.services
}

output "run_command" {
  description = "systemctl command to enable and start all containers"
  value       = local.run_command
}
