output "ec2_private_ip" {
  value       = try(module.ec2[0].private_ip, "")
  description = "GitLab Runner manager private IP."
}
