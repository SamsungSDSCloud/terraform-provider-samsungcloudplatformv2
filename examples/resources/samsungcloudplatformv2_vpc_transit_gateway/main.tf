provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_transitgateway" "tgw01" {
  name          = var.name
  tags          = var.tags
  description   = var.description
}

