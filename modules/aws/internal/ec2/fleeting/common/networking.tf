##############
# Networking #
##############

resource "aws_vpc" "jobs-vpc" {
  cidr_block = var.aws_vpc_cidr

  tags = merge(var.labels, {
    Name = "${var.name}"
  })
}

# adding tags to aws_route_table causes an error: query returned no results
data "aws_route_table" "jobs-route-table" {
  vpc_id = aws_vpc.jobs-vpc.id

  depends_on = [
    aws_vpc.jobs-vpc
  ]
}

resource "aws_internet_gateway" "internet-access" {
  vpc_id = aws_vpc.jobs-vpc.id

  tags = merge(var.labels, {
    Name = "${var.name}"
  })
}

# both `name` and `tags` are unsupported arguments
resource "aws_route" "internet-route" {
  route_table_id         = data.aws_route_table.jobs-route-table.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.internet-access.id
}

resource "aws_subnet" "jobs-vpc-subnet" {
  vpc_id     = aws_vpc.jobs-vpc.id
  cidr_block = var.asg_subnet_cidr

  availability_zone = var.aws_zone

  map_public_ip_on_launch = true

  tags = merge(var.labels, {
    Name = "${var.name}"
  })
}

resource "aws_security_group" "jobs-security-group" {
  name   = "${var.name}"
  vpc_id = aws_vpc.jobs-vpc.id

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
    Name = "${var.name}"
  })
}

