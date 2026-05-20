provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_vpc_subnet_vip" "subnet_vip_detail" {
  subnet_id = var.subnet_id
  vip_id    = var.vip_id
}
