#######################
# METADATA VALIDATION #
#######################

module "validate_name" {
  source = "../../internal/validation/name"
  name   = var.metadata.name
}

module "validate_support" {
  source   = "../../internal/validation/support"
  use_case = "any"
  use_case_support = tomap({
    "any" = "experimental"
  })
  min_support = var.metadata.min_support
}

###################
# VPC PROD MODULE #
###################

resource "aws_vpc" "vpc" {
  cidr_block           = var.cidr
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = merge(var.metadata.labels, {
    Name = var.metadata.name
  })
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id

  tags = merge(var.metadata.labels, {
    Name = var.metadata.name
  })
}

resource "aws_subnet" "jobs_vpc_subnet" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = var.subnet_cidr

  availability_zone = var.zone

  map_public_ip_on_launch = true

  tags = merge(var.metadata.labels, {
    Name = var.metadata.name
  })
}

# adding tags to aws_route_table causes an error: query returned no results
resource "aws_route_table" "rtb_public" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_route_table_association" "rta_subnet_public" {
  subnet_id      = aws_subnet.jobs_vpc_subnet.id
  route_table_id = aws_route_table.rtb_public.id
}

# both `name` and `tags` are unsupported arguments
resource "aws_route" "internet_route" {
  route_table_id         = aws_route_table.rtb_public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.igw.id
}

