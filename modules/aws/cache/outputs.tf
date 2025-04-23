output "enabled" {
  value = tobool(true)
}

output "server_address" {
  description = "Address of the S3 bucket server"
  value       = tostring("s3.amazonaws.com")
}

output "bucket_name" {
  description = "Name of the created bucket"
  value       = tostring(aws_s3_bucket.cache.id)
}

output "bucket_location" {
  description = "AWS region of the cache bucket"
  value       = tostring(aws_s3_bucket.cache.region)
}

output "access_key_id" {
  description = "Access key ID for the user with access to the cache bucket"
  value       = tostring(aws_iam_access_key.cache_bucket_user_key.id)
}

output "secret_access_key" {
  description = "Secret access key for the user with access to the cache bucket"
  value       = tostring(aws_iam_access_key.cache_bucket_user_key.secret)
  sensitive   = true
}
