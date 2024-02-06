output "bucket_name" {
  value       = module.cache.bucket_name
  description = "Name of the GCS bucket created for storing runner remote cache"
}
