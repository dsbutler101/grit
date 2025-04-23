output "enabled" {
  value = tobool(true)
}

output "instance_group_name" {
  value       = tostring(try(module.gce[0].instance_group_name, ""))
  description = "Name of the created instance group (when 'gce' fleeting service in use)"
}
