##############################
# Runner manager credentials #
##############################

locals {
  create = var.create_key_pair != null
}

# both `name` and `tags` are unsupported arguments
resource "tls_private_key" "aws_runner_key_pair" {
  count = local.create ? 1 : 0

  algorithm = var.create_key_pair.algorithm
  rsa_bits  = var.create_key_pair.size
}

resource "aws_key_pair" "aws_runner_key_pair" {
  count = local.create ? 1 : 0

  key_name   = var.name
  public_key = tls_private_key.aws_runner_key_pair[0].public_key_openssh

  tags = merge(var.labels, {
    name = var.name
  })
}
