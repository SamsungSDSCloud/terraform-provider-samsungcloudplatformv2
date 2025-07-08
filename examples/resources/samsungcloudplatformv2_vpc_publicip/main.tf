provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_publicip" "publicip" {
    description = var.publicip_description
    type = var.publicip_type
}
