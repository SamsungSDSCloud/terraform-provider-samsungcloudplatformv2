provider "samsungcloudplatformv2" {
}

data "samsungcloudplatformv2_firewall_firewall_rule" "firewall_rule" {
  id = var.id
  firewall_id = var.firewall_id
  dst_ip = var.dst_ip
}
