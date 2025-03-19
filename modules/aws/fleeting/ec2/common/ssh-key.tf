#################################
# Autoscaling group credentials #
#################################

# both `name` and `tags` are unsupported arguments
resource "tls_private_key" "aws_jobs_private_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "jobs_key_pair" {
  key_name   = var.name
  public_key = tls_private_key.aws_jobs_private_key.public_key_openssh

  tags = merge(var.labels, {
    Name = var.name
  })
}
