provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_transit_gateways" "vpctransitgateway" {
  size = var.size
  sort = var.sort
  page = var.page
  state = var.state
}
