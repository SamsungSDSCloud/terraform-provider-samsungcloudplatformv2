provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_firewall_firewall_rule" "firewallrule" {
  firewall_id = var.firewall_id
  firewall_rule_create = var.firewall_rule
}
