output "fleeting" {
  description = "Security group for the ephemeral fleeting VMs"
  value       = aws_security_group.jobs_security_group
}

output "runner_manager" {
  description = "Security group for the runner manager"
  value       = aws_security_group.manager_sg
}
