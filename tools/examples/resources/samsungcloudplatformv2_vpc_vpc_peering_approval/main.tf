provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_vpc_peering_approval" "my_approve" {
  vpc_peering_id = var.vpc_peering_id
  type           = var.type
}
