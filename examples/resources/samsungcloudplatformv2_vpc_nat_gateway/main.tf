provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_nat_gateway" "natgateway" {
    subnet_id = var.natgateway_subnet_id
    publicip_ids = var.natgateway_publicip_ids
    multi_zone_enabled = var.natgateway_multi_zone_enabled
    description = var.natgateway_description
}
