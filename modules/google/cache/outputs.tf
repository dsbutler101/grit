output "bucket_name" {
  value       = google_storage_bucket.cache-bucket.name
  description = "Name of the GCS bucket created for storing runner remote cache"
}
