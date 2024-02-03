output "fleeting_access_key_id" {
  description = "The non-secret ID of the service account access key"
  value       = module.iam.fleeting_access_key_id
}

output "fleeting_secret_access_key" {
  description = "The secret access key of the service account"
  value       = module.iam.fleeting_secret_access_key
  sensitive   = true
}
