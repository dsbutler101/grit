output "enabled" {
  value = module.cache.enabled
}

output "server_address" {
  description = "Address of the S3 bucket server"
  value       = module.cache.server_address
}

output "bucket_name" {
  description = "Name of the created bucket"
  value       = module.cache.bucket_name
}

output "bucket_location" {
  description = "AWS region of the cache bucket"
  value       = module.cache.bucket_location
}

output "access_key_id" {
  description = "Access key ID for the user with access to the cache bucket"
  value       = module.cache.access_key_id
}

output "secret_access_key" {
  description = "Secret access key for the user with access to the cache bucket"
  value       = module.cache.secret_access_key
  sensitive   = true
}
