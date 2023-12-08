output "ssh_key_pem" {
  value = module.instance_group.ssh_key_pem
}

output "fleeting_service_account_access_key_id" {
  value = module.instance_group.fleeting_service_account_access_key_id
}

output "fleeting_service_account_secret_access_key" {
  value = module.instance_group.fleeting_service_account_secret_access_key
}

