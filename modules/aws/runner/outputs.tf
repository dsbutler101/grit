output "enabled" {
  value = tobool(true)
}

output "ec2_private_ip" {
  value       = tostring(try(module.ec2[0].private_ip, ""))
  description = "GitLab Runner manager private IP."
}

output "ec2_public_ip" {
  value       = tostring(try(module.ec2[0].public_ip, ""))
  description = "GitLab Runner manager public IP."
}

output "ec2_instance_name" {
  value       = tostring(try(module.ec2[0].instance_name, ""))
  description = "GitLab Runner manager instance name."
}

output "ec2_instance_id" {
  value       = tostring(try(module.ec2[0].instance_id, ""))
  description = "GitLab Runner manager instance ID."
}

output "ec2_runner_wrapper_socket_path" {
  value       = tostring(try(module.ec2[0].runner_wrapper_socket_path, ""))
  description = "The address of the runner wrapper on the manager"
}

output "ec2_ssh_key_openssh_pem" {
  value       = tostring(try(module.ec2[0].ssh_key_openssh_pem, ""))
  description = "GitLab Runner manager SSH key."
  sensitive   = true
}
