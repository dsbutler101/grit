output "id" {
  value = google_compute_network.default.id
}

output "subnetwork_ids" {
  value = tomap({
    for name, subnet in google_compute_subnetwork.subnetwork : name => subnet.id
  })
}

output "subnetwork_cidrs" {
  value = tomap({
    for name, subnet in google_compute_subnetwork.subnetwork : name => subnet.ip_cidr_range
  })
}