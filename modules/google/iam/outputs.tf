output "service_account_email" {
  value       = google_service_account.default.email
  description = "Email of the created service account"
}
