output "jobs_host_resource_group_outputs" {
  description = "The outputs of a jobs host resource group for hosting MacOS VMs"
  value       = aws_cloudformation_stack.jobs_cloudformation_stack.outputs
}

output "license_config_arn" {
  description = "The ARN for licenses used to deploy MacOS instances"
  value       = aws_licensemanager_license_configuration.license_config.arn
}
