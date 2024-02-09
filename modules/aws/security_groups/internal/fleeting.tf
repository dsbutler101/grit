resource "aws_security_group" "jobs_security_group" {
  name   = "${var.name} fleet"
  vpc_id = var.vpc_id

  ingress {
    description = "SSH from outside"
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "ICMP from remote"
    protocol    = "icmp"
    from_port   = 8
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.labels, {
    Name = "${var.name} fleet"
  })
}