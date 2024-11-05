output "external_ip" {
  value = google_compute_instance.prometheus-server.network_interface[0].access_config[0].nat_ip
}

output "internal_hostname" {
  value = "${google_compute_instance.prometheus-server.name}.c.${var.google_project}.internal"
}
