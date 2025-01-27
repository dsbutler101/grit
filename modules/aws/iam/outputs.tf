output "fleeting_access_key_id" {
  description = "The non-secret ID of the service account access key"
  value       = aws_iam_access_key.fleeting-service-account-key.id
}

output "fleeting_secret_access_key" {
  description = "The secret access key of the service account"
  value       = aws_iam_access_key.fleeting-service-account-key.secret
  sensitive   = true
}
