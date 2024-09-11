output "ami_id" {
  value = local.manifest[local.use_case][var.region].id
}
