output "ssh_key_pem" {
  description = "The pem file with SSH key for access to the autoscaling group instances"
  value       = try(module.ec2[0].ssh_key_pem, "")
  sensitive   = true
}

output "ssh_key_pem_name" {
  description = "The name of the pem file with SSH key for access to the autoscaling group instances"
  value       = try(module.ec2[0].ssh_key_pem_name, "")
}

output "autoscaling_group_name" {
  description = "The name of the autoscaling group created"
  value       = try(module.ec2[0].autoscaling_group_name, "")
}
