output "ssh_key_pem" {
  value       = module.dev.ssh_key_pem
  description = "The pem file with SSH key for access to the autoscaling group instances"
  sensitive   = true
}

output "fleeting_service_account_access_key_id" {
  value       = module.dev.fleeting_service_account_access_key_id
  description = "The access key ID for access to the fleeting service account"
  sensitive   = true
}

output "fleeting_service_account_secret_access_key" {
  value       = module.dev.fleeting_service_account_secret_access_key
  description = "The secret access key for access to the fleeting service account"
  sensitive   = true
}
