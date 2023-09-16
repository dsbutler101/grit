output "output_map" {
  description = "Outputs from GCE resources"
  value = tomap({
    "linux"   = module.linux[0].output_map,
    "windows" = try(module.windows[0].output_map, null)
  })
}