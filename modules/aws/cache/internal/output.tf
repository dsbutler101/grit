output "enabled" {
  value = true
}

output "server_address" {
  value = "s3.amazonaws.com"
}

output "bucket_name" {
  value = aws_s3_bucket.cache.id
}

output "bucket_location" {
  value = aws_s3_bucket.cache.region
}

output "access_key_id" {
  value = aws_iam_access_key.cache-bucket-user-key.id
}

output "secret_access_key" {
  value     = aws_iam_access_key.cache-bucket-user-key.secret
  sensitive = true
}
