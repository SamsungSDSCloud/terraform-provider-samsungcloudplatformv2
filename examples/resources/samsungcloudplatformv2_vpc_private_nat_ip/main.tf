provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_private_nat_ip" "privatenatip" {
  private_nat_id = var.private_nat_id
  ip_address = var.private_nat_ip_ip_address
  description = var.private_nat_ip_description
}
