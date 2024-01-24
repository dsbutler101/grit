###################################
# SSH key for accessing instances #
###################################

output "ssh_key_pem_name" {
  value       = try(module.ec2[0].ssh_key_pem_name, "")
  description = "The pem file with SSH key for access to the autoscaling group instances"
  sensitive   = true
}


output "ssh_key_pem" {
  value       = try(module.ec2[0].ssh_key_pem, "")
  description = "The pem file with SSH key for access to the autoscaling group instances"
  sensitive   = true
}

#####################
# Autoscaling Group #
#####################

output "autoscaling_group_name" {
  value = try(module.ec2[0].autoscaling_group_name, "")
}
