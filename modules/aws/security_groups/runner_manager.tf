resource "aws_security_group" "manager_sg" {
  name   = "${var.metadata.name} manager"
  vpc_id = var.vpc_id

  dynamic "ingress" {
    for_each = var.manager_inbound_sg_rules
    content {
      from_port   = ingress.value.from_port
      to_port     = ingress.value.to_port
      protocol    = ingress.value.protocol
      cidr_blocks = ingress.value.cidr_blocks
      description = try(ingress.value.description, "Managed by Terraform")
    }
  }

  dynamic "egress" {
    for_each = var.manager_outbound_sg_rules
    content {
      from_port   = egress.value.from_port
      to_port     = egress.value.to_port
      protocol    = egress.value.protocol
      cidr_blocks = egress.value.cidr_blocks
      description = try(egress.value.description, "Managed by Terraform")
    }
  }

  tags = merge(module.labels.merged, {
    name = "${var.metadata.name} manager"
  })
}
