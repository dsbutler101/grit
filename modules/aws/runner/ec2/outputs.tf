output "private_ip" {
  value = aws_instance.runner-manager.private_ip
}

output "deprecated_warning" {
  value       = var.vpc.subnet_id != null ? "Warning: The 'subnet_id' variable is deprecated. Please use 'subnet_ids' instead." : null
  description = "A warning for providind subnet_id"
}

output "public_ip" {
  value = aws_instance.runner-manager.public_ip
}
