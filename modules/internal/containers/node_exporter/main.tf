locals {
  container_config = {
    name    = var.service_name
    image   = "${var.registry}/${var.image_path}:${var.image_tag}"
    network = "host"
    pid     = "host"
    volumes = ["/:/host:ro,rslave"]
    command = <<EOT
    --web.listen-address=0.0.0.0:${var.port} \
    --path.rootfs=/host
    EOT
    service_options = [{
      ExecStartPost  = "/sbin/iptables -A INPUT -p tcp -m tcp --dport ${var.port} -j ACCEPT"
      TimeoutStopSec = "30"
    }]
  }
}
