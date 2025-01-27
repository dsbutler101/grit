resource "aws_security_group" "jobs_security_group" {
  name   = "${var.metadata.name} fleet"
  vpc_id = var.vpc_id

  dynamic "ingress" {
    for_each = var.fleeting_inbound_sg_rules
    content {
      from_port   = ingress.value.from_port
      to_port     = ingress.value.to_port
      protocol    = ingress.value.protocol
      cidr_blocks = ingress.value.cidr_blocks
      description = try(ingress.value.description, "Managed by Terraform")
    }
  }

  dynamic "egress" {
    for_each = var.fleeting_outbound_sg_rules
    content {
      from_port   = egress.value.from_port
      to_port     = egress.value.to_port
      protocol    = egress.value.protocol
      cidr_blocks = egress.value.cidr_blocks
      description = try(egress.value.description, "Managed by Terraform")
    }
  }

  tags = merge(var.metadata.labels, {
    Name = "${var.metadata.name} fleet"
  })
}
