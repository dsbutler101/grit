output "private_ip" {
  value = aws_instance.runner_manager.private_ip
}

output "deprecated_warning" {
  value       = var.vpc.subnet_id != null ? "Warning: The 'subnet_id' variable is deprecated. Please use 'subnet_ids' instead." : null
  description = "A warning for providind subnet_id"
}

output "public_ip" {
  value = aws_instance.runner_manager.public_ip
}

output "grit_runner_manager" {
  value = var.runner_wrapper.enabled ? {
    instance_name   = aws_instance.runner_manager.id
    address         = aws_instance.runner_manager.private_ip
    wrapper_address = var.runner_wrapper.socket_path
  } : {}
}
