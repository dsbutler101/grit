output "internal_ip" {
  description = "Internal IP of the deployer runner manager"
  value       = google_compute_instance.runner_manager.network_interface[0].network_ip
}

output "external_ip" {
  description = "External IP of deployed runner manager only if access config is enabled"
  value       = var.access_config_enabled ? google_compute_instance.runner_manager.network_interface[0].access_config[0].nat_ip : null
}

output "internal_hostname" {
  description = "Internal hostname of the runner manager instance"
  value       = "${google_compute_instance.runner_manager.name}.c.${var.google_project}.internal"
}

output "instance_name" {
  description = "Name of the runner manager instance"
  value       = google_compute_instance.runner_manager.name
}

output "instance_id" {
  description = "ID of the runner manager instance"
  value       = google_compute_instance.runner_manager.instance_id
}

output "wrapper_address" {
  description = "Address of the process wrapper socket on the runner manager instance (if enabled)"
  value       = local.runner_wrapper.enabled ? local.runner_wrapper.socket_path : ""
}
