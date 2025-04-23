output "enabled" {
  value = tobool(true)
}

output "bucket_name" {
  value       = tostring(google_storage_bucket.cache_bucket.name)
  description = "Name of the GCS bucket created for storing runner remote cache"
}
