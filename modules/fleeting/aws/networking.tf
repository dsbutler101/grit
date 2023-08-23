##############
# Networking #
##############

resource "aws_vpc" "jobs-vpc" {
  cidr_block = var.aws_vpc_cidr

  tags = merge(local.tags, {
    Name = "jobs-vpc"
  })
}

data "aws_route_table" "jobs-vpc" {
  vpc_id = aws_vpc.jobs-vpc.id

  depends_on = [
    aws_vpc.jobs-vpc
  ]
}

resource "aws_internet_gateway" "internet-access" {
  vpc_id = aws_vpc.jobs-vpc.id

  tags = local.tags
}

resource "aws_route" "internet-access" {
  route_table_id         = data.aws_route_table.jobs-vpc.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.internet-access.id
}

resource "aws_subnet" "jobs-vpc-subnet" {
  for_each = var.autoscaling_groups

  vpc_id     = aws_vpc.jobs-vpc.id
  cidr_block = each.value.subnet_cidr

  availability_zone = var.aws_zone

  map_public_ip_on_launch = true

  tags = merge(local.tags, {
    Name = each.key
  })
}

resource "aws_security_group" "jobs-security-group" {
  name   = var.shard
  vpc_id = aws_vpc.jobs-vpc.id

  ingress {
    description = "SSH from runner manager"
    protocol    = "tcp"
    from_port   = 22
    to_port     = 22
    cidr_blocks = [local.gcp_runner_manager_subnetwork_cidr]
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

  tags = local.tags
}