output "enabled" {
  value = tobool(true)
}

output "ami_id" {
  value = tostring(local.manifest[local.use_case][var.region].id)
}
