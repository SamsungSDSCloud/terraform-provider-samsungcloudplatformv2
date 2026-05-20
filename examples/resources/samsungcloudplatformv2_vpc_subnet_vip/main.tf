provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_vpc_subnet_vip" "my_subnet_vip" {
  subnet_id          = var.subnet_id
  virtual_ip_address = var.virtual_ip_address
  description        = var.description
}
