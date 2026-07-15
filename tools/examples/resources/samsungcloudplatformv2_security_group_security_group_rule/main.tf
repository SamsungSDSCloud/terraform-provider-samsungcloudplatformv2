provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_security_group_security_group_rule" "securitygrouprule" {
  security_group_id = var.security_group_id
  ethertype = var.ethertype
  protocol = var.protocol
  port_range_min = var.port_range_min
  port_range_max = var.port_range_max
  remote_ip_prefix = var.remote_ip_prefix
  description = var.description
  direction = var.direction
}
