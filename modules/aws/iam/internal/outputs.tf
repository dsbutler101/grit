output "fleeting_access_key_id" {
  description = "TODO"
  value       = aws_iam_access_key.fleeting-service-account-key.id
}

output "fleeting_secret_access_key" {
  description = "TODO"
  value       = aws_iam_access_key.fleeting-service-account-key.secret
  sensitive   = true
}
