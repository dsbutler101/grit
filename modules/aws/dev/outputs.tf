output "ssh_key_pem" {
  value     = try(module.dev-module.ssh_key_pem, "")
  sensitive = true
}

output "fleeting_service_account_access_key_id" {
  value     = try(module.dev-module.fleeting_service_account_access_key_id, "")
  sensitive = true
}

output "fleeting_service_account_secret_access_key" {
  value     = try(module.dev-module.fleeting_service_account_secret_access_key, "")
  sensitive = true
}

