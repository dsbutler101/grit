output "external_ip" {
  value = google_compute_instance.runner-manager.network_interface[0].access_config[0].nat_ip
}

output "internal_hostname" {
  value = "${google_compute_instance.runner-manager.name}.c.${var.google_project}.internal"
}
