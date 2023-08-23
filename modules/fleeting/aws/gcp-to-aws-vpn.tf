locals {
  gcp_asn = 64514
  aws_asn = 64515
}

############
# GCP side #
############

resource "google_compute_router" "gcp-to-aws" {
  name    = "${var.shard}-gcp-to-aws"
  region  = var.gcp_region
  network = var.gcp_runner_manager_vpc_link

  bgp {
    asn = local.gcp_asn

    advertise_mode    = "CUSTOM"
    advertised_groups = ["ALL_SUBNETS"]
  }
}

resource "google_compute_ha_vpn_gateway" "gcp-to-aws" {
  name    = "${var.shard}-gcp-to-aws"
  region  = var.gcp_region
  network = var.gcp_runner_manager_vpc_link
}

resource "google_compute_external_vpn_gateway" "gcp-to-aws" {
  name            = "${var.shard}-gcp-to-aws"
  redundancy_type = "FOUR_IPS_REDUNDANCY"

  interface {
    id         = 0
    ip_address = aws_vpn_connection.aws-to-gcp-1.tunnel1_address
  }

  interface {
    id         = 1
    ip_address = aws_vpn_connection.aws-to-gcp-1.tunnel2_address
  }

  interface {
    id         = 2
    ip_address = aws_vpn_connection.aws-to-gcp-2.tunnel1_address
  }

  interface {
    id         = 3
    ip_address = aws_vpn_connection.aws-to-gcp-2.tunnel2_address
  }
}

resource "google_compute_vpn_tunnel" "gcp-to-aws-1" {
  name                            = "${var.shard}-gcp-to-aws-1"
  region                          = var.gcp_region
  vpn_gateway                     = google_compute_ha_vpn_gateway.gcp-to-aws.id
  vpn_gateway_interface           = 0
  peer_external_gateway           = google_compute_external_vpn_gateway.gcp-to-aws.id
  peer_external_gateway_interface = 0
  router                          = google_compute_router.gcp-to-aws.id
  shared_secret                   = aws_vpn_connection.aws-to-gcp-1.tunnel1_preshared_key
}

resource "google_compute_vpn_tunnel" "gcp-to-aws-2" {
  name                            = "${var.shard}-gcp-to-aws-2"
  region                          = var.gcp_region
  vpn_gateway                     = google_compute_ha_vpn_gateway.gcp-to-aws.id
  vpn_gateway_interface           = 0
  peer_external_gateway           = google_compute_external_vpn_gateway.gcp-to-aws.id
  peer_external_gateway_interface = 1
  router                          = google_compute_router.gcp-to-aws.id
  shared_secret                   = aws_vpn_connection.aws-to-gcp-1.tunnel2_preshared_key
}

resource "google_compute_vpn_tunnel" "gcp-to-aws-3" {
  name                            = "${var.shard}-gcp-to-aws-3"
  region                          = var.gcp_region
  vpn_gateway                     = google_compute_ha_vpn_gateway.gcp-to-aws.id
  vpn_gateway_interface           = 1
  peer_external_gateway           = google_compute_external_vpn_gateway.gcp-to-aws.id
  peer_external_gateway_interface = 2
  router                          = google_compute_router.gcp-to-aws.id
  shared_secret                   = aws_vpn_connection.aws-to-gcp-2.tunnel1_preshared_key
}

resource "google_compute_vpn_tunnel" "gcp-to-aws-4" {
  name                            = "${var.shard}-gcp-to-aws-4"
  region                          = var.gcp_region
  vpn_gateway                     = google_compute_ha_vpn_gateway.gcp-to-aws.id
  vpn_gateway_interface           = 1
  peer_external_gateway           = google_compute_external_vpn_gateway.gcp-to-aws.id
  peer_external_gateway_interface = 3
  router                          = google_compute_router.gcp-to-aws.id
  shared_secret                   = aws_vpn_connection.aws-to-gcp-2.tunnel2_preshared_key
}

resource "google_compute_router_interface" "gcp-to-aws-1" {
  name       = "${var.shard}-gcp-to-aws-1"
  region     = var.gcp_region
  router     = google_compute_router.gcp-to-aws.name
  ip_range   = "${aws_vpn_connection.aws-to-gcp-1.tunnel1_cgw_inside_address}/30"
  vpn_tunnel = google_compute_vpn_tunnel.gcp-to-aws-1.name
}

resource "google_compute_router_interface" "gcp-to-aws-2" {
  name       = "${var.shard}-gcp-to-aws-2"
  region     = var.gcp_region
  router     = google_compute_router.gcp-to-aws.name
  ip_range   = "${aws_vpn_connection.aws-to-gcp-1.tunnel2_cgw_inside_address}/30"
  vpn_tunnel = google_compute_vpn_tunnel.gcp-to-aws-2.name
}

resource "google_compute_router_interface" "gcp-to-aws-3" {
  name       = "${var.shard}-gcp-to-aws-3"
  region     = var.gcp_region
  router     = google_compute_router.gcp-to-aws.name
  ip_range   = "${aws_vpn_connection.aws-to-gcp-2.tunnel1_cgw_inside_address}/30"
  vpn_tunnel = google_compute_vpn_tunnel.gcp-to-aws-3.name
}

resource "google_compute_router_interface" "gcp-to-aws-4" {
  name       = "${var.shard}-gcp-to-aws-4"
  region     = var.gcp_region
  router     = google_compute_router.gcp-to-aws.name
  ip_range   = "${aws_vpn_connection.aws-to-gcp-2.tunnel2_cgw_inside_address}/30"
  vpn_tunnel = google_compute_vpn_tunnel.gcp-to-aws-4.name
}

resource "google_compute_router_peer" "gcp-to-aws-1" {
  name            = "${var.shard}-gcp-to-aws-1"
  region          = var.gcp_region
  router          = google_compute_router.gcp-to-aws.name
  interface       = "${var.shard}-gcp-to-aws-1"
  peer_ip_address = aws_vpn_connection.aws-to-gcp-1.tunnel1_vgw_inside_address
  peer_asn        = local.aws_asn
}

resource "google_compute_router_peer" "gcp-to-aws-2" {
  name            = "${var.shard}-gcp-to-aws-2"
  region          = var.gcp_region
  router          = google_compute_router.gcp-to-aws.name
  interface       = "${var.shard}-gcp-to-aws-2"
  peer_ip_address = aws_vpn_connection.aws-to-gcp-1.tunnel2_vgw_inside_address
  peer_asn        = local.aws_asn
}

resource "google_compute_router_peer" "gcp-to-aws-3" {
  name            = "${var.shard}-gcp-to-aws-3"
  region          = var.gcp_region
  router          = google_compute_router.gcp-to-aws.name
  interface       = "${var.shard}-gcp-to-aws-3"
  peer_ip_address = aws_vpn_connection.aws-to-gcp-2.tunnel1_vgw_inside_address
  peer_asn        = local.aws_asn
}

resource "google_compute_router_peer" "gcp-to-aws-4" {
  name            = "${var.shard}-gcp-to-aws-4"
  region          = var.gcp_region
  router          = google_compute_router.gcp-to-aws.name
  interface       = "${var.shard}-gcp-to-aws-4"
  peer_ip_address = aws_vpn_connection.aws-to-gcp-2.tunnel2_vgw_inside_address
  peer_asn        = local.aws_asn
}

############
# AWS side #
############

resource "aws_customer_gateway" "aws-to-gcp-1" {
  type = "ipsec.1"

  bgp_asn    = local.gcp_asn
  ip_address = google_compute_ha_vpn_gateway.gcp-to-aws.vpn_interfaces[0].ip_address

  tags = local.tags
}

resource "aws_customer_gateway" "aws-to-gcp-2" {
  type = "ipsec.1"

  bgp_asn    = local.gcp_asn
  ip_address = google_compute_ha_vpn_gateway.gcp-to-aws.vpn_interfaces[1].ip_address

  tags = local.tags
}

resource "aws_vpn_gateway" "aws-to-gcp" {
  amazon_side_asn = local.aws_asn

  vpc_id = aws_vpc.jobs-vpc.id

  tags = local.tags
}

resource "aws_vpn_connection" "aws-to-gcp-1" {
  type                = aws_customer_gateway.aws-to-gcp-1.type
  customer_gateway_id = aws_customer_gateway.aws-to-gcp-1.id
  vpn_gateway_id      = aws_vpn_gateway.aws-to-gcp.id

  tunnel1_phase1_encryption_algorithms = ["AES256"]
  tunnel1_phase2_encryption_algorithms = ["AES256"]
  tunnel1_phase1_integrity_algorithms  = ["SHA2-256"]
  tunnel1_phase2_integrity_algorithms  = ["SHA2-256"]
  tunnel1_phase1_dh_group_numbers      = [14]
  tunnel1_phase2_dh_group_numbers      = [14]

  tunnel2_phase1_encryption_algorithms = ["AES256"]
  tunnel2_phase2_encryption_algorithms = ["AES256"]
  tunnel2_phase1_integrity_algorithms  = ["SHA2-256"]
  tunnel2_phase2_integrity_algorithms  = ["SHA2-256"]
  tunnel2_phase1_dh_group_numbers      = [14]
  tunnel2_phase2_dh_group_numbers      = [14]

  tags = local.tags
}

resource "aws_vpn_connection" "aws-to-gcp-2" {
  type                = aws_customer_gateway.aws-to-gcp-2.type
  customer_gateway_id = aws_customer_gateway.aws-to-gcp-2.id
  vpn_gateway_id      = aws_vpn_gateway.aws-to-gcp.id

  tunnel1_phase1_encryption_algorithms = ["AES256"]
  tunnel1_phase2_encryption_algorithms = ["AES256"]
  tunnel1_phase1_integrity_algorithms  = ["SHA2-256"]
  tunnel1_phase2_integrity_algorithms  = ["SHA2-256"]
  tunnel1_phase1_dh_group_numbers      = [14]
  tunnel1_phase2_dh_group_numbers      = [14]

  tunnel2_phase1_encryption_algorithms = ["AES256"]
  tunnel2_phase2_encryption_algorithms = ["AES256"]
  tunnel2_phase1_integrity_algorithms  = ["SHA2-256"]
  tunnel2_phase2_integrity_algorithms  = ["SHA2-256"]
  tunnel2_phase1_dh_group_numbers      = [14]
  tunnel2_phase2_dh_group_numbers      = [14]

  tags = local.tags
}

resource "aws_vpn_gateway_route_propagation" "aws-to-gcp" {
  vpn_gateway_id = aws_vpn_gateway.aws-to-gcp.id
  route_table_id = data.aws_route_table.jobs-vpc.id
}
