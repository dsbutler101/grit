output "enabled" {
  value = tobool(true)
}

output "id" {
  description = "ID of the created VPC"
  value       = tostring(google_compute_network.default.id)
}

output "subnetwork_ids" {
  description = "Map of created subnetwork IDs"
  value = tomap({
    for name, subnet in google_compute_subnetwork.subnetwork : name => subnet.id
  })
}

output "subnetwork_cidrs" {
  description = "Map of created subnetwork CIDRs"
  value = tomap({
    for name, subnet in google_compute_subnetwork.subnetwork : name => subnet.ip_cidr_range
  })
}

output "nat_ips" {
  description = "External IP addresses allocated to NAT"
  value       = google_compute_router_nat.nat.nat_ips
}
