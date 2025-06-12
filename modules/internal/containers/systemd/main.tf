
locals {
  services = [for c in var.containers : {
    path        = "${var.service_path}/${c.name}.service"
    owner       = var.owner
    permissions = var.permissions
    content     = templatefile("${path.module}/systemd-service.tftpl", c)
  }]

  service_names = [for c in var.containers : "${c.name}.service"]
  run_command   = length(var.containers) == 0 ? "" : "systemctl daemon-reload && systemctl enable --now ${join(" ", local.service_names)}"
}
