provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_subnet_vip_nat_ip" "my_subnet_vip_nat_ip" {
  subnet_id   = var.subnet_id
  vip_id      = var.vip_id
  nat_type    = var.nat_type
  publicip_id = var.publicip_id
}
