output "output_map" {
  description = "Outputs from EC2 resources"
  value = tomap({
    "macos" = module.macos[0].output_map,
  })
}