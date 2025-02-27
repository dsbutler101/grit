locals {
  node_exporter_dir            = "/opt/node_exporter/"
  node_exporter_installer_path = "${local.node_exporter_dir}/install_node_exporter.sh"

  write_files_config = [{
    path        = local.node_exporter_installer_path
    content     = local.install_node_exporter_script
    owner       = "root:root"
    permissions = "0755"
  }]

  install_node_exporter_commands = [
    local.node_exporter_installer_path,
  ]

  install_node_exporter_script = templatefile("${path.module}/templates/install_node_exporter.sh", {
    node_exporter_dir     = local.node_exporter_dir
    node_exporter_version = var.node_exporter_version
    node_exporter_port    = var.node_exporter_port
  })
}
