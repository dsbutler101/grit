output "gke-cluster" {
  value = try(module.gke-cluster[0].output_map, null)
}