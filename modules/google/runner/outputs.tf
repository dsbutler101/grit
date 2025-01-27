output "external_ip" {
  description = "External IP of deployed runner manager"
  value       = google_compute_instance.runner-manager.network_interface[0].access_config[0].nat_ip
}

output "internal_hostname" {
  description = "Internal hostname of the runner manager instance"
  value       = "${google_compute_instance.runner-manager.name}.c.${var.google_project}.internal"
}
