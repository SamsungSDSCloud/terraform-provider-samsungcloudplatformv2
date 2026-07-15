provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_transit_gateway" "vpctransitgateway" {
  id = var.id
}
