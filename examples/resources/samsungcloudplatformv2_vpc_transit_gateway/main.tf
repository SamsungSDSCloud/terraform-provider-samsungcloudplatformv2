provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_transit_gateway" "tgw01" {
  name          = var.name
  tags          = var.tags
  description   = var.description
}

