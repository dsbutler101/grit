output "output_map" {
  description = "Outputs from the Fleeting Instance Group"
  value = tomap({
    "ssh_key_pem" = module.instance_group.ssh_key_pem,
    "fleeting_service_account_access_key_id" = module.instance_group.fleeting_service_account_access_key_id,
    "fleeting_service_account_secret_access_key" = module.instance_group.fleeting_service_account_secret_access_key,
//    "autoscaling_group_names" = module.instance_group.autoscaling_group_names,
  })
}