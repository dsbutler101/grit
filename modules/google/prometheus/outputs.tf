output "enabled" {
  value = tobool(true)
}

output "external_ip" {
  description = "External IP of the Prometheus server instance"
  value       = tostring(google_compute_instance.prometheus_server.network_interface[0].access_config[0].nat_ip)
}

output "internal_hostname" {
  description = "Internal hostname of the prometheus server instance"
  value       = tostring("${google_compute_instance.prometheus_server.name}.c.${var.google_project}.internal")
}
