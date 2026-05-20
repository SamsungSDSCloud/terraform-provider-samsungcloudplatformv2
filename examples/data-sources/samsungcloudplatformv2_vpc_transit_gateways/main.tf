provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_transit_gateways" "vpctransitgateway" {
  firewall_connection_state = var.firewall_connection_state
  name                      = var.name
  id                        = var.id
  size                      = var.size
  sort                      = var.sort
  page                      = var.page
  state                     = var.state
}
