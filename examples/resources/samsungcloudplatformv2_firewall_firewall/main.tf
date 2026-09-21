provider "samsungcloudplatformv2" {
}

resource "samsungcloudplatformv2_firewall_firewall" "firewall" {
  flavor_name = var.flavor_name
  loggable = var.loggable
}
