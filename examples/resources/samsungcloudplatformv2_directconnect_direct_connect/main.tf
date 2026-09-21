provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_directconnect_direct_connect" "directconnect" {
  name = var.dcon_name
  vpc_id = var.dcon_vpc_id
  bandwidth = var.dcon_bandwidth
  description = var.dcon_description
  uplink_active_zone = var.dcon_uplink_active_zone
  uplink_standby_zone = var.dcon_uplink_standby_zone
  firewall_enabled = var.dcon_firewall_enabled
  firewall_loggable = var.dcon_firewall_loggable
  tags = var.tags
}