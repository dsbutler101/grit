output "ssh_key_pem_name" {
  value     = module.common.ssh_key_pem_name
  sensitive = true
}


output "ssh_key_pem" {
  value     = module.common.ssh_key_pem
  sensitive = true
}

output "autoscaling_group_name" {
  value = module.common.autoscaling_group_name
}

output "deprecated_warning" {
  value       = var.vpc.subnet_id != null ? "Warning: The 'subnet_id' variable is deprecated. Please use 'subnet_ids' instead." : null
  description = "A warning for providind subnet_id"
}
