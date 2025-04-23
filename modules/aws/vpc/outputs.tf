output "enabled" {
  value = tobool(true)
}

output "id" {
  value = tostring(aws_vpc.vpc.id)
}

output "subnet_ids" {
  value = tolist([aws_subnet.jobs_vpc_subnet.id])
}
