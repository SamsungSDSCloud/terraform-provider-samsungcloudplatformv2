provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_firewall_firewall_rule" "firewallrule1" {
  firewall_id = var.firewall_id
  firewall_rule_create = var.firewall_rule1
}

resource "samsungcloudplatformv2_firewall_firewall_rule" "firewallrule2" {
  firewall_id = var.firewall_id
  firewall_rule_create = var.firewall_rule2
  depends_on = [samsungcloudplatformv2_firewall_firewall_rule.firewallrule1]
}

resource "samsungcloudplatformv2_firewall_firewall_rule" "firewallrule3" {
  firewall_id = var.firewall_id
  firewall_rule_create = var.firewall_rule3
  depends_on = [samsungcloudplatformv2_firewall_firewall_rule.firewallrule2]
}
