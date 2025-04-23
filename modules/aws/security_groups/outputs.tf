output "enabled" {
  value = tobool(true)
}

output "fleeting_id" {
  description = "Security group for the ephemeral fleeting VMs"
  value       = tostring(aws_security_group.jobs_security_group.id)
}

output "runner_manager_id" {
  description = "Security group for the runner manager"
  value       = tostring(aws_security_group.manager_sg.id)
}
