output "jobs-host-resource-group-outputs" {
  description = "The outputs of a jobs host resource group for hosting MacOS VMs"
  value       = aws_cloudformation_stack.jobs-cloudformation-stack.outputs
}

output "license-config-arn" {
  description = "The ARN for licenses used to deploy MacOS instances"
  value       = aws_licensemanager_license_configuration.license-config.arn
}

