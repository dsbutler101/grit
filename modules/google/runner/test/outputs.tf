output "external_ip" {
  value       = module.runner.external_ip
  description = "External IP of deployed runner manager"
}
