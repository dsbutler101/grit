output "write_files_config" {
  value       = local.write_files_config
  description = "List of cloud-init write_file directives needed to install the node exporter"
}

output "commands" {
  value       = local.install_node_exporter_commands
  description = "List of cloud-init commands needed to install the node exporter"
}
