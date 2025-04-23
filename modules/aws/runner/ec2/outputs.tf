output "private_ip" {
  value = aws_instance.runner_manager.private_ip
}

output "public_ip" {
  value = aws_instance.runner_manager.public_ip
}

output "instance_name" {
  value = local.instance_name
}

output "instance_id" {
  value = aws_instance.runner_manager.id
}

output "runner_wrapper_socket_path" {
  value = var.runner_wrapper.socket_path
}

output "ssh_key_openssh_pem" {
  value     = try(tls_private_key.aws_runner_key_pair[0].private_key_openssh, "")
  sensitive = true
}
