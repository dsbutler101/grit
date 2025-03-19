output "external_ip" {
  description = "External IP of deployed runner manager only if access config is enabled"
  value       = var.access_config_enabled ? google_compute_instance.runner_manager.network_interface[0].access_config[0].nat_ip : null
}

output "internal_hostname" {
  description = "Internal hostname of the runner manager instance"
  value       = "${google_compute_instance.runner_manager.name}.c.${var.google_project}.internal"
}
