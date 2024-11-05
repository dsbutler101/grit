output "external_ip" {
  description = "External IP of the Prometheus server instance"
  value       = module.prometheus.external_ip
}

output "internal_hostname" {
  description = "Internal hostname of the prometheus server instance"
  value       = module.prometheus.internal_hostname
}