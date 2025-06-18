# output "runner_manager_external_ip" {
#   value = module.runner.external_ip
# }

output "nat_ips" {
  description = "External IP addresses allocated to NAT"
  value       = module.vpc.nat_ips
}

output "prometheus_external_ip" {
  value = try(module.prometheus[0].external_ip, "")
}

output "prometheus_internal_hostname" {
  value = try(module.prometheus[0].internal_hostname, "")
}
