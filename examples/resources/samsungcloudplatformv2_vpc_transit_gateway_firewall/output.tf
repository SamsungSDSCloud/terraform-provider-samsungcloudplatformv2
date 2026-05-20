output "my_tgw_firewall" {
  value = {
    transit_gateway : samsungcloudplatformv2_vpc_transit_gateway_firewall.my_tgw_firewall.transit_gateway
  }
}
