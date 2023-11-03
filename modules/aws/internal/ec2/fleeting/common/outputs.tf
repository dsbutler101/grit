###################################
# SSH key for accessing instances #
###################################

output "ssh_key_pem_name" {
  value       = aws_key_pair.jobs.key_name
  description = "The pem file with SSH key for access to the autoscaling group instances"
  sensitive   = true
}

output "ssh_key_pem" {
  value       = tls_private_key.aws-jobs.private_key_pem
  description = "The pem file with SSH key for access to the autoscaling group instances"
  sensitive   = true
}

###############################
# Service account credentials #
###############################

output "fleeting_service_account_access_key_id" {
  value       = aws_iam_access_key.fleeting-service-account-key.id
  description = "The access key ID for access to the fleeting service account"
  sensitive   = true
}

output "fleeting_service_account_secret_access_key" {
  value       = aws_iam_access_key.fleeting-service-account-key.secret
  description = "The secret access key for access to the fleeting service account"
  sensitive   = true
}

#####################
# Autoscaling Group #
#####################

output "autoscaling_group_names" {
  value = [
    for asg in aws_autoscaling_group.fleeting-asg : asg.name
  ]
}
