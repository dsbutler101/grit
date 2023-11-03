output "ssh_key_pem" {
  value = module.ec2[0].ssh_key_pem
}

output "fleeting_service_account_access_key_id" {
  value = module.ec2[0].fleeting_service_account_access_key_id
}

output "fleeting_service_account_secret_access_key" {
  value = module.ec2[0].fleeting_service_account_secret_access_key
}