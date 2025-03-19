output "id" {
  value = aws_vpc.vpc.id
}

output "subnet_id" {
  value = aws_subnet.jobs_vpc_subnet.id
}

output "subnet_ids" {
  value = [aws_subnet.jobs_vpc_subnet.id]
}
