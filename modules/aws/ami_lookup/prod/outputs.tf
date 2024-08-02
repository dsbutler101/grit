output "ami_id" {
  value = local.manifest[var.use_case][var.region].id
}
