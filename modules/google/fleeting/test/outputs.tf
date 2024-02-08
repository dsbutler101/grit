output "instance_group_name" {
  value       = try(module.gce.instance_group_name, "")
  description = "Name of the created instance group (when 'gce' fleeting service in use)"
}
