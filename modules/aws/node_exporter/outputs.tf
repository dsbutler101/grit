output "enabled" {
  value = tobool(true)
}

output "write_files_config" {
  value       = tolist(local.write_files_config)
  description = "List of cloud-init write_file directives needed to install the node exporter"
}

output "commands" {
  value       = tolist(local.install_node_exporter_commands)
  description = "List of cloud-init commands needed to install the node exporter"
}

output "port" {
  value       = tonumber(var.node_exporter_port)
  description = "Port that the node exporter is listening on"
}

output "version" {
  value       = tostring(var.node_exporter_version)
  description = "The version of the Prometheus node exporter to install"
}
