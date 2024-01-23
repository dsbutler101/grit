output "id" {
  value       = module.vpc.id
  description = "ID of the created VPC"
}

output "subnetwork_ids" {
  value       = module.vpc.subnetwork_ids
  description = "Map of created subnetwork IDs"
}

output "subnetwork_cidrs" {
  value       = module.vpc.subnetwork_cidrs
  description = "Map of created subnetwork CIDRs"
}
