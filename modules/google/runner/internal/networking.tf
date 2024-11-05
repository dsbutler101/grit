resource "google_compute_firewall" "runner-manager-ssh-access" {
  name    = "${var.name}-runner-manager-ssh-access"
  network = var.vpc.id

  direction = "INGRESS"
  priority  = 1000

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]

  target_tags = [local.runner_manager_tag]
}

resource "google_compute_firewall" "additional-rules" {
  for_each = var.runner_manager_additional_firewall_rules

  name    = "${var.name}-runner-${each.key}-access"
  network = var.vpc.id

  direction = each.value.direction
  priority  = each.value.priority

  dynamic "allow" {
    for_each = each.value.allow

    content {
      protocol = allow.value.protocol
      ports    = [for port in allow.value.ports : tostring(port)]
    }
  }

  dynamic "deny" {
    for_each = each.value.deny

    content {
      protocol = deny.value.protocol
      ports    = [for port in deny.value.ports : tostring(port)]
    }
  }

  source_ranges = each.value.source_ranges

  target_tags = [local.runner_manager_tag]
}