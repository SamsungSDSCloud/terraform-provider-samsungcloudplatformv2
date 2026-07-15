provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_subnet_vip_port" "my_subnet_vip_port" {
  subnet_id = var.subnet_id
  vip_id    = var.vip_id
  port_id   = var.port_id
}
