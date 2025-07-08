provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_nat_gateway" "natgateway" {
    subnet_id = var.natgateway_subnet_id
    publicip_id = var.natgateway_publicip_id
    description = var.natgateway_description
}
