output "external_ip" {
  value       = module.runner.external_ip
  description = "External IP of deployed runner manager"
}

output "internal_hostname" {
  description = "Internal hostname of the runner manager instance"
  value       = module.runner.internal_hostname
}
