###################################
# SSH key for accessing instances #
###################################

output "ssh_key_pem_name" {
  value       = module.ec2-instance-group[0].ssh_key_pem_name
  description = "The pem file with SSH key for access to the autoscaling group instances"
  sensitive   = true
}


output "ssh_key_pem" {
  value       = module.ec2-instance-group[0].ssh_key_pem
  description = "The pem file with SSH key for access to the autoscaling group instances"
  sensitive   = true
}

###############################
# Service account credentials #
###############################

output "fleeting_service_account_access_key_id" {
  value       = module.ec2-instance-group[0].fleeting_service_account_access_key_id
  description = "The access key ID for access to the fleeting service account"
  sensitive   = true
}

output "fleeting_service_account_secret_access_key" {
  value       = module.ec2-instance-group[0].fleeting_service_account_secret_access_key
  description = "The secret access key for access to the fleeting service account"
  sensitive   = true
}

#####################
# Autoscaling Group #
#####################

output "autoscaling_group_names" {
  value = module.ec2-instance-group[0].autoscaling_group_name
}

