variable "labels" {
  type = map(any)
}

variable "name" {
  type = string
}

variable "aws_zone" {
  type = string
}

variable "aws_vpc_cidr" {
  type = string
}

variable "aws_vpc_subnet_cidr" {
  type = string
}