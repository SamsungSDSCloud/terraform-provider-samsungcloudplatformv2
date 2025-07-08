provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_directconnect_direct_connect" "directconnect" {
  name = var.dcon_name
  vpc_id = var.dcon_vpc_id
  bandwidth = var.dcon_bandwidth
  description = var.dcon_description
}
