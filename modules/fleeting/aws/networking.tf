##############
# Networking #
##############

resource "aws_vpc" "jobs-vpc" {
  cidr_block = var.aws_vpc_cidr

  tags = merge(var.labels, {
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

  tags = var.labels
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

  tags = merge(var.labels, {
    Name = each.key
  })
}
