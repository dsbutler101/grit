###################################
# SSH key for accessing instances #
###################################

output "ssh_key_pem_name" {
  value       = aws_key_pair.jobs_key_pair.key_name
  description = "The pem file with SSH key for access to the autoscaling group instances"
  sensitive   = true
}

output "ssh_key_pem" {
  value       = tls_private_key.aws_jobs_private_key.private_key_pem
  description = "The pem file with SSH key for access to the autoscaling group instances"
  sensitive   = true
}

#####################
# Autoscaling Group #
#####################

output "autoscaling_group_name" {
  value = aws_autoscaling_group.fleeting_asg.name
}
