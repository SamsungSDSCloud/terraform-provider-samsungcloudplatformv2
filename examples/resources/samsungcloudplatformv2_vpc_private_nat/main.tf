provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_private_nat" "privatenat" {
  cidr = var.private_nat_cidr
  description = var.private_nat_description
  name = var.private_nat_name
  service_resource_id = var.private_nat_service_resource_id
  service_type = var.private_nat_service_type
  tags = var.private_nat_tags
}
