output "fleeting_access_key_id" {
  description = "The non-secret ID of the service account access key"
  value       = aws_iam_access_key.fleeting_service_account_key.id
}

output "fleeting_secret_access_key" {
  description = "The secret access key of the service account"
  value       = aws_iam_access_key.fleeting_service_account_key.secret
  sensitive   = true
}
