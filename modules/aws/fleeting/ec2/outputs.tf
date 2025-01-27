output "ssh_key_pem_name" {
  value     = module.common.ssh_key_pem_name
  sensitive = true
}


output "ssh_key_pem" {
  value     = module.common.ssh_key_pem
  sensitive = true
}

output "autoscaling_group_name" {
  value = module.common.autoscaling_group_name
}
