output "external_ip" {
  description = "External IP of the Prometheus server instance"
  value       = google_compute_instance.prometheus-server.network_interface[0].access_config[0].nat_ip
}

output "internal_hostname" {
  description = "Internal hostname of the prometheus server instance"
  value       = "${google_compute_instance.prometheus-server.name}.c.${var.google_project}.internal"
}