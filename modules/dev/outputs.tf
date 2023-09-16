output "dev" {
  description = "Outputs from the dev module"
  value = tomap({
    "ec2" = try(module.ec2[0].output_map, null),
    "gce" = try(module.gce[0].output_map, null),
  })
}