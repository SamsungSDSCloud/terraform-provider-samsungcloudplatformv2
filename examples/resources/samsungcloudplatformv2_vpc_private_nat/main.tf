provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_private_nat" "privatenat" {
  name = var.private_nat_name
  direct_connect_id = var.private_nat_direct_connect_id
  cidr = var.private_nat_cidr
  description = var.private_nat_description
  tags = var.private_nat_tags
}
