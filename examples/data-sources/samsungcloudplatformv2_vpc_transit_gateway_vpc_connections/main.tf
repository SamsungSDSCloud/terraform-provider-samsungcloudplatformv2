provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_transit_gateway_vpc_connections" "transitgateway_vpc_connections" {
  size = var.size
  sort = var.sort
  page = var.page
  id = var.id
  transit_gateway_id = var.transit_gateway_id
  vpc_id = var.vpc_id
  vpc_name = var.vpc_name
  state = var.state
}
