
locals {
  services = [for c in var.containers : {
    file_name    = "${c.name}.service"
    file_content = templatefile("${path.module}/systemd-service.tftpl", c)
  }]

  run_command = length(var.containers) == 0 ? "" : "systemctl daemon-reload && systemctl enable --now ${join(" ", local.services[*].file_name)}"
}
