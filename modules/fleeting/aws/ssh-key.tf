#################################
# Autoscaling group credentials #
#################################

resource "tls_private_key" "aws-jobs" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "jobs" {
  key_name   = "${var.shard}-ssh-key"
  public_key = tls_private_key.aws-jobs.public_key_openssh

  tags = local.tags
}
